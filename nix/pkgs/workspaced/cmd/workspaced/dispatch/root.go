package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"
	"workspaced/cmd/workspaced/dispatch/apply"
	"workspaced/cmd/workspaced/dispatch/audio"
	"workspaced/cmd/workspaced/dispatch/backup"
	"workspaced/cmd/workspaced/dispatch/brightness"
	"workspaced/cmd/workspaced/dispatch/config"
	"workspaced/cmd/workspaced/dispatch/doctor"
	"workspaced/cmd/workspaced/dispatch/history"
	"workspaced/cmd/workspaced/dispatch/nix"
	"workspaced/cmd/workspaced/dispatch/notification"
	"workspaced/cmd/workspaced/dispatch/open"
	"workspaced/cmd/workspaced/dispatch/palette"
	"workspaced/cmd/workspaced/dispatch/plan"
	"workspaced/cmd/workspaced/dispatch/power"
	"workspaced/cmd/workspaced/dispatch/screen"
	"workspaced/cmd/workspaced/dispatch/shell"
	"workspaced/cmd/workspaced/dispatch/sudo"
	"workspaced/cmd/workspaced/dispatch/sync"
	"workspaced/cmd/workspaced/dispatch/template"
	"workspaced/cmd/workspaced/dispatch/wallpaper"
	"workspaced/cmd/workspaced/system/screenshot"
	"workspaced/cmd/workspaced/system/workspace"
	libconfig "workspaced/pkg/config"
	"workspaced/pkg/exec"
	"workspaced/pkg/types"

	_ "workspaced/pkg/driver/prelude"

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
		// Load config to initialize driver weights and other global states
		if _, err := libconfig.Load(); err != nil {
			slog.Debug("failed to load config in dispatch", "error", err)
		}

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

		if remoteCmd == "sync" || remoteCmd == "config" || remoteCmd == "is" || remoteCmd == "shell" {
			return nil
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
	cmd.AddCommand(doctor.Command)
	cmd.AddCommand(history.GetCommand())

	cmd.AddCommand(nix.GetCommand())
	cmd.AddCommand(notification.GetCommand())
	cmd.AddCommand(open.GetCommand())
	cmd.AddCommand(palette.GetCommand())
	cmd.AddCommand(plan.GetCommand())
	cmd.AddCommand(power.GetCommand())
	cmd.AddCommand(screen.GetCommand())
	cmd.AddCommand(screenshot.GetCommand())
	cmd.AddCommand(shell.GetCommand())
	cmd.AddCommand(sudo.GetCommand())
	cmd.AddCommand(sync.GetCommand())
	cmd.AddCommand(template.GetCommand())
	cmd.AddCommand(wallpaper.GetCommand())
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
			return net.DialTimeout("unix", socketPath, 1*time.Second)
		},
	}

	conn, _, err := dialer.Dial("ws://localhost/ws", nil)
	if err != nil {
		slog.Info("daemon not reachable, running locally", "error", err)
		return "", false, nil
	}
	defer func() { _ = conn.Close() }()

	// Get client binary hash
	clientHash, _ := exec.GetBinaryHash()

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
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Debug("ws read error", "error", err)
			}
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
				// Check if daemon is restarting itself
				if resp.Error == "DAEMON_RESTARTING" || resp.Error == "DAEMON_RESTART_NEEDED" {
					slog.Info("daemon restarting with new binary, retrying locally")

					// Daemon is exec'ing itself, just wait a bit and run locally
					// Next command will connect to the new daemon
					time.Sleep(200 * time.Millisecond)

					// Run locally this time, next call will hit new daemon
					return "", false, nil
				}
				return "", true, fmt.Errorf("%s", resp.Error)
			}
			return "", true, nil
		}
	}
}
