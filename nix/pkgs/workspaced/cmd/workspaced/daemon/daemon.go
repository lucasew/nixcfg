package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/dispatch"
)

var Command = &cobra.Command{
	Use:   "daemon",
	Short: "Run the workspaced daemon",
	Run: func(c *cobra.Command, args []string) {
		if err := RunDaemon(); err != nil {
			slog.Error("daemon failure", "error", err)
			os.Exit(1)
		}
	},
}

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "workspaced.sock")
}

func RunDaemon() error {
	var listener net.Listener

	listeners, err := activation.Listeners()
	if err == nil && len(listeners) > 0 {
		listener = listeners[0]
	} else {
		socketPath := getSocketPath()
		os.Remove(socketPath)
		l, err := net.Listen("unix", socketPath)
		if err != nil {
			return fmt.Errorf("failed to listen on socket: %w", err)
		}
		listener = l
	}
	defer listener.Close()

	slog.Info("listening", "address", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("accept error", "error", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req dispatch.Request
	if err := decoder.Decode(&req); err != nil {
		slog.Warn("invalid request", "error", err)
		encoder.Encode(dispatch.Response{Error: fmt.Sprintf("invalid request: %v", err)})
		return
	}

	slog.Info("executing command", "command", req.Command, "args", req.Args)

	output, err := ExecuteViaCobra(req)
	resp := dispatch.Response{Output: output}
	if err != nil {
		slog.Error("command failed", "command", req.Command, "error", err)
		resp.Error = err.Error()
	}

	encoder.Encode(resp)
}

func ExecuteViaCobra(req dispatch.Request) (string, error) {
	fullArgs := append([]string{req.Command}, req.Args...)

	root := dispatch.Command
	buf := new(bytes.Buffer)

	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(fullArgs)

	// Inject WORKSPACED_DAEMON to prevent recursion in child processes
	env := append(req.Env, "WORKSPACED_DAEMON=1")

	// Inject daemon_mode flag and environment
	ctx := context.WithValue(context.Background(), "env", env)
	ctx = context.WithValue(ctx, "daemon_mode", true)

	err := root.ExecuteContext(ctx)
	return buf.String(), err
}
