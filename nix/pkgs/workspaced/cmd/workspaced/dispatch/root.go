package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"workspaced/cmd/workspaced/dispatch/apply"
	"workspaced/cmd/workspaced/dispatch/audio"
	"workspaced/cmd/workspaced/dispatch/backup"
	"workspaced/cmd/workspaced/dispatch/brightness"
	"workspaced/cmd/workspaced/dispatch/config"
	"workspaced/cmd/workspaced/dispatch/demo"
	"workspaced/cmd/workspaced/dispatch/history"
	"workspaced/cmd/workspaced/dispatch/media"
	"workspaced/cmd/workspaced/dispatch/menu"
	"workspaced/cmd/workspaced/dispatch/nix"
	"workspaced/cmd/workspaced/dispatch/notification"
	"workspaced/cmd/workspaced/dispatch/power"
	"workspaced/cmd/workspaced/dispatch/screen"
	"workspaced/cmd/workspaced/dispatch/screenshot"
	"workspaced/cmd/workspaced/dispatch/setup"
	"workspaced/cmd/workspaced/dispatch/sudo"
	"workspaced/cmd/workspaced/dispatch/wallpaper"
	"workspaced/cmd/workspaced/dispatch/webapp"
	"workspaced/cmd/workspaced/dispatch/workspace"
	"workspaced/cmd/workspaced/is"
	"workspaced/pkg/common"
	"workspaced/pkg/types"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "dispatch",
		Short:            "Dispatch workspace commands",
		TraverseChildren: true,
	}

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		isDaemon := false

		val := ctx.Value(types.DaemonModeKey)
		if os.Getenv("WORKSPACED_DAEMON") == "1" {
			isDaemon = true
		}
		if val == true {
			isDaemon = true
		}

		if isDaemon {
			return nil
		}

		var remoteCmd string
		var remoteArgs []string

		for i, arg := range os.Args {
			if arg == "dispatch" && i+1 < len(os.Args) {
				remoteCmd = os.Args[i+1]
				remoteArgs = os.Args[i+2:]
				break
			}
		}

		if remoteCmd == "" {
			return nil
		}

		if remoteCmd == "sudo" && len(remoteArgs) > 0 {
			switch remoteArgs[0] {
			case "approve", "reject", "add":
				return nil
			}
		}

		if remoteCmd == "history" && len(remoteArgs) > 0 {
			switch remoteArgs[0] {
			case "search", "list", "ingest":
				return nil
			}
		}

		output, connected, err := TryRemoteRaw(remoteCmd, remoteArgs)
		if connected {
			if output != "" {
				fmt.Print(output)
			}
			if err != nil {
				return err
			}
			os.Exit(0)
		}

		return nil
	}

	cmd.AddCommand(apply.GetCommand())
	cmd.AddCommand(audio.GetCommand())
	cmd.AddCommand(backup.GetCommand())
	cmd.AddCommand(brightness.GetCommand())
	cmd.AddCommand(config.GetCommand())
	cmd.AddCommand(demo.GetCommand())
	cmd.AddCommand(history.GetCommand())
	cmd.AddCommand(is.GetCommand())
	cmd.AddCommand(media.GetCommand())
	cmd.AddCommand(menu.GetCommand())
	cmd.AddCommand(nix.GetCommand())
	cmd.AddCommand(notification.GetCommand())
	cmd.AddCommand(power.GetCommand())
	cmd.AddCommand(screen.GetCommand())
	cmd.AddCommand(screenshot.GetCommand())
	cmd.AddCommand(setup.GetCommand())
	cmd.AddCommand(sudo.GetCommand())
	cmd.AddCommand(wallpaper.GetCommand())
	cmd.AddCommand(webapp.GetCommand())
	cmd.AddCommand(workspace.GetCommand())

	return cmd
}

func FindCommand(name string, args []string) (*cobra.Command, []string, error) {
	return NewCommand().Find(append([]string{name}, args...))
}

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "workspaced.sock")
}

func TryRemoteRaw(cmdName string, args []string) (string, bool, error) {
	socketPath := getSocketPath()
	slog.Info("connecting to daemon", "socket", socketPath, "cmd", cmdName, "args", args)

	dialer := websocket.Dialer{
		NetDial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, 5*time.Second)
		},
	}

	conn, _, err := dialer.Dial("ws://localhost/ws", nil)
	if err != nil {
		slog.Info("daemon not reachable, running locally", "error", err)
		return "", false, nil
	}
	defer func() { _ = conn.Close() }()

	// Get client binary hash
	clientHash, _ := common.GetBinaryHash()

	req := types.Request{
		Command:    cmdName,
		Args:       args,
		Env:        os.Environ(),
		BinaryHash: clientHash,
	}

	// Send request as a StreamPacket
	payload, _ := json.Marshal(req)
	packet := types.StreamPacket{
		Type:    "request",
		Payload: payload,
	}

	if err := conn.WriteJSON(packet); err != nil {
		return "", true, fmt.Errorf("failed to send request: %w", err)
	}

	for {
		var packet types.StreamPacket
		if err := conn.ReadJSON(&packet); err != nil {
			return "", true, fmt.Errorf("failed to read response: %w", err)
		}

		switch packet.Type {
		case "log":
			var entry types.LogEntry
			if err := json.Unmarshal(packet.Payload, &entry); err != nil {
				continue
			}
			level := slog.LevelInfo
			switch entry.Level {
			case "DEBUG":
				level = slog.LevelDebug
			case "WARN":
				level = slog.LevelWarn
			case "ERROR":
				level = slog.LevelError
			}
			attrs := []any{}
			for k, v := range entry.Attrs {
				attrs = append(attrs, slog.Any(k, v))
			}
			slog.Log(context.Background(), level, entry.Message, attrs...)
		case "stdout":
			var out string
			if err := json.Unmarshal(packet.Payload, &out); err == nil {
				fmt.Print(out)
			}
		case "stderr":
			var out string
			if err := json.Unmarshal(packet.Payload, &out); err == nil {
				fmt.Fprint(os.Stderr, out)
			}
		case "result":
			var resp types.Response
			if err := json.Unmarshal(packet.Payload, &resp); err != nil {
				return "", true, fmt.Errorf("failed to parse result: %w", err)
			}
			if resp.Error != "" {
				// Check if daemon needs restart
				if resp.Error == "DAEMON_RESTART_NEEDED" {
					slog.Info("daemon binary outdated, restarting daemon and retrying")

					// Restart the daemon via systemd
					cmd := exec.Command("systemctl", "--user", "restart", "workspaced.service")
					_ = cmd.Run()

					// Wait a bit for daemon to restart
					time.Sleep(500 * time.Millisecond)

					// Retry the command locally (new daemon will be picked up next time)
					// Return not connected so caller runs locally, which will connect to new daemon
					return "", false, nil
				}
				return "", true, fmt.Errorf("%s", resp.Error)
			}
			return "", true, nil
		}
	}
}
