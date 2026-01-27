package dispatch

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	Command.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		// Detect if we are running INSIDE the daemon to avoid recursion/loops.
		if os.Getenv("WORKSPACED_DAEMON") == "1" {
			slog.Info("running inside daemon, skipping remote execution")
			return nil
		}

		ctx := c.Context()
		if ctx != nil && ctx.Value("daemon_mode") == true {
			return nil // We are the daemon, proceed to local logic
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
	slog.Info("connecting to daemon", "socket", socketPath, "command", c.Name())

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		slog.Info("daemon not reachable, running locally", "error", err)
		return "", false, nil
	}
	defer conn.Close()

	req := Request{
		Command: c.Name(),
		Args:    args,
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
