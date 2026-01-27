package dispatch

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	Command.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		isDaemon := false

		val := ctx.Value(DaemonModeKey)
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

		// We are the client. Try to connect to daemon.
		output, connected, err := TryRemote(c, args)
		if connected {
			if output != "" {
				fmt.Print(output)
			}
			if err != nil {
				return err
			}
			os.Exit(0) // Clean exit after remote execution
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

func TryRemote(c *cobra.Command, args []string) (string, bool, error) {
	socketPath := getSocketPath()
	fullPath := GetFullCommandPath(c)
	if len(fullPath) == 0 {
		return "", false, nil
	}
	slog.Info("connecting to daemon", "socket", socketPath, "path", fullPath, "args", args)

	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		slog.Info("daemon not reachable, running locally", "error", err)
		return "", false, nil
	}
	defer conn.Close()

	// Set a read/write deadline to avoid hanging forever
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	req := Request{
		Command: fullPath[0],
		Args:    append(fullPath[1:], args...),
		Env:     os.Environ(),
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		return "", true, fmt.Errorf("failed to send request: %w", err)
	}

	var resp Response
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&resp); err != nil {
		return "", true, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.Error != "" {
		return resp.Output, true, fmt.Errorf(resp.Error)
	}

	return resp.Output, true, nil
}
