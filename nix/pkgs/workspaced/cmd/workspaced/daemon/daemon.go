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
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/pkg/drivers/screen"
	"workspaced/pkg/types"
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

	// Start background tasks
	go monitorCapsLock()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("accept error", "error", err)
			continue
		}
		slog.Info("accepted connection", "remote", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func monitorCapsLock() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Find capslock files
	matches, _ := filepath.Glob("/sys/class/leds/*capslock/brightness")
	if len(matches) == 0 {
		slog.Warn("no capslock leds found for monitoring")
		return
	}

	lastActive := false

	for range ticker.C {
		active := false
		for _, m := range matches {
			data, err := os.ReadFile(m)
			if err == nil && strings.TrimSpace(string(data)) == "1" {
				active = true
				break
			}
		}

		if active && !lastActive {
			slog.Info("capslock activated, turning off screen")
			ctx := context.Background()
			if err := screen.SetDPMS(ctx, false); err != nil {
				slog.Error("failed to turn off screen from capslock", "error", err)
			}
		}
		lastActive = active
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req types.Request
	if err := decoder.Decode(&req); err != nil {
		slog.Warn("failed to decode request", "error", err)
		encoder.Encode(types.Response{Error: fmt.Sprintf("invalid request: %v", err)})
		return
	}

	slog.Info("executing command", "command", req.Command, "args", req.Args)

	output, err := ExecuteViaCobra(req)
	resp := types.Response{Output: output}
	if err != nil {
		slog.Error("command failed", "command", req.Command, "args", req.Args, "error", err)
		resp.Error = err.Error()
	}

	slog.Info("sending response", "output_len", len(output))
	if err := encoder.Encode(resp); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func ExecuteViaCobra(req types.Request) (string, error) {
	targetCmd, targetArgs, err := dispatch.FindCommand(req.Command, req.Args)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	targetCmd.SetOut(buf)
	targetCmd.SetErr(buf)
	targetCmd.SetArgs(targetArgs)

	// Prepare context
	env := append(req.Env, "WORKSPACED_DAEMON=1")
	ctx := context.WithValue(context.Background(), types.EnvKey, env)
	ctx = context.WithValue(ctx, types.DaemonModeKey, true)

	// Set context on the command
	targetCmd.SetContext(ctx)

	// Manually run PersistentPreRunE of the PARENTS from root down to target
	var parents []*cobra.Command
	for curr := targetCmd; curr != nil; curr = curr.Parent() {
		parents = append([]*cobra.Command{curr}, parents...)
	}

	for _, p := range parents {
		if p.PersistentPreRunE != nil {
			if err := p.PersistentPreRunE(targetCmd, targetArgs); err != nil {
				return buf.String(), err
			}
		} else if p.PersistentPreRun != nil {
			p.PersistentPreRun(targetCmd, targetArgs)
		}
	}

	// Run the command
	if targetCmd.RunE != nil {
		err = targetCmd.RunE(targetCmd, targetArgs)
	} else if targetCmd.Run != nil {
		targetCmd.Run(targetCmd, targetArgs)
	} else {
		err = fmt.Errorf("command has no run implementation")
	}

	return buf.String(), err
}
