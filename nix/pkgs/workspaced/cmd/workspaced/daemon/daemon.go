package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"io"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/pkg/common"
	"workspaced/pkg/types"
)

type StreamPacketWriter struct {
	Out  chan<- types.StreamPacket
	Type string
}

func (w *StreamPacketWriter) Write(p []byte) (n int, err error) {
	payload, _ := json.Marshal(string(p))
	w.Out <- types.StreamPacket{
		Type:    w.Type,
		Payload: payload,
	}
	return len(p), nil
}

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

	server := &http.Server{
		Handler: http.HandlerFunc(handleWS),
	}

	return server.Serve(listener)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade error", "error", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Output channel
	outCh := make(chan types.StreamPacket, 1000)
	done := make(chan struct{})

	// Pump goroutine: channel -> websocket
	go func() {
		defer close(done)
		for packet := range outCh {
			if err := conn.WriteJSON(packet); err != nil {
				slog.Error("ws write error", "error", err)
				cancel()
				return
			}
		}
	}()

	// Read request
	var req types.Request
	if err := conn.ReadJSON(&req); err != nil {
		slog.Warn("ws read request error", "error", err)
		return
	}

	slog.Info("executing command", "command", req.Command, "args", req.Args)

	// Create logger
	handler := &common.ChannelLogHandler{
		Out:    outCh,
		Parent: slog.Default().Handler(),
		Ctx:    ctx,
	}
	logger := slog.New(handler)

	// Build context
	stdout := &StreamPacketWriter{Out: outCh, Type: "stdout"}
	stderr := &StreamPacketWriter{Out: outCh, Type: "stderr"}

	ctx = context.WithValue(ctx, types.LoggerKey, logger)
	ctx = context.WithValue(ctx, types.StdoutKey, stdout)
	ctx = context.WithValue(ctx, types.StderrKey, stderr)
	env := append(req.Env, "WORKSPACED_DAEMON=1")
	ctx = context.WithValue(ctx, types.EnvKey, env)
	ctx = context.WithValue(ctx, types.DaemonModeKey, true)

	output, err := ExecuteViaCobra(ctx, req, stdout, stderr)

	resp := types.Response{Output: output}
	if err != nil {
		slog.Error("command failed", "command", req.Command, "error", err)
		resp.Error = err.Error()
	}

	payload, _ := json.Marshal(resp)
	outCh <- types.StreamPacket{
		Type:    "result",
		Payload: payload,
	}

	close(outCh)
	<-done
}

func ExecuteViaCobra(ctx context.Context, req types.Request, stdout, stderr io.Writer) (string, error) {
	targetCmd, targetArgs, err := dispatch.FindCommand(req.Command, req.Args)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	targetCmd.SetOut(io.MultiWriter(buf, stdout))
	targetCmd.SetErr(io.MultiWriter(buf, stderr))
	targetCmd.SetArgs(targetArgs)
	targetCmd.SetContext(ctx)

	if err := targetCmd.ParseFlags(targetArgs); err != nil {
		return buf.String(), err
	}
	argList := targetCmd.Flags().Args()

	var parents []*cobra.Command
	for curr := targetCmd; curr != nil; curr = curr.Parent() {
		parents = append([]*cobra.Command{curr}, parents...)
	}

	for _, p := range parents {
		if p.PersistentPreRunE != nil {
			if err := p.PersistentPreRunE(targetCmd, argList); err != nil {
				return buf.String(), err
			}
		} else if p.PersistentPreRun != nil {
			p.PersistentPreRun(targetCmd, argList)
		}
	}

	if targetCmd.RunE != nil {
		err = targetCmd.RunE(targetCmd, argList)
	} else if targetCmd.Run != nil {
		targetCmd.Run(targetCmd, argList)
	} else {
		err = fmt.Errorf("command has no run implementation")
	}

	return buf.String(), err
}
