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
	"workspaced/pkg/types"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func init() {
	Command.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		isDaemon := false

		val := ctx.Value(types.DaemonModeKey)
		slog.Debug("checking daemon mode", "command", c.Name(), "ctx_val", val, "env_var", os.Getenv("WORKSPACED_DAEMON"))

		if os.Getenv("WORKSPACED_DAEMON") == "1" {
			isDaemon = true
		}
		if val == true {
			isDaemon = true
		}

		if isDaemon {
			slog.Info("running inside daemon, skipping remote execution", "command", c.Name())
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
	defer conn.Close()

	req := types.Request{
		Command: cmdName,
		Args:    args,
		Env:     os.Environ(),
	}

	if err := conn.WriteJSON(req); err != nil {
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
		case "result":
			var resp types.Response
			if err := json.Unmarshal(packet.Payload, &resp); err != nil {
				return "", true, fmt.Errorf("failed to parse result: %w", err)
			}
			if resp.Error != "" {
				return resp.Output, true, fmt.Errorf(resp.Error)
			}
			return resp.Output, true, nil
		}
	}
}
