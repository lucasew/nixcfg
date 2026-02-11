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
	"syscall"
	"time"

	"image"
	_ "image/png"
	"io"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/pkg/db"
	"workspaced/pkg/exec"
	"workspaced/pkg/icons"
	"workspaced/pkg/logging"
	"workspaced/pkg/media"
	"workspaced/pkg/tray"
	_ "workspaced/pkg/tray/dbus"
	"workspaced/pkg/types"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	shouldRestartDaemon bool
	initialMtime        time.Time
)

func init() {
	var err error
	initialMtime, err = exec.GetBinaryMtime()
	if err != nil {
		slog.Warn("failed to get initial binary mtime", "error", err)
	}
}

func HasBinaryChanged() bool {
	if initialMtime.IsZero() {
		return false
	}
	currentMtime, err := exec.GetBinaryMtime()
	return err == nil && !currentMtime.Equal(initialMtime)
}

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
		try, _ := c.Flags().GetBool("try")
		if try {
			socketPath := getSocketPath()
			conn, err := net.DialTimeout("unix", socketPath, 200*time.Millisecond)
			if err == nil {
				conn.Close()
				slog.Info("daemon already running, exiting")
				os.Exit(0)
			}
		}

		if err := RunDaemon(); err != nil && err != http.ErrServerClosed {
			slog.Error("daemon failure", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	Command.Flags().Bool("try", false, "Exit if daemon is already running")
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

	// Ensure logs go to stderr
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("daemon starting", "pid", os.Getpid())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.Open()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer database.Close()

	go media.Watch(ctx)

	// Initialize tray
	go func() {
		t, err := tray.GetDefault()
		if err != nil {
			slog.Debug("no tray driver found, skipping", "error", err)
			return
		}

		iconPath, err := icons.GetIconPath(ctx, "https://github.com")
		var icon image.Image
		if err == nil {
			f, err := os.Open(iconPath)
			if err == nil {
				defer f.Close()
				icon, _, _ = image.Decode(f)
			}
		}

		t.SetState(tray.State{
			Title: "workspaced",
			Icon:  icon,
			Menu: []tray.MenuItem{
				{
					Label: "Apply",
					Callback: func() {
						slog.Info("tray: triggering apply")
						// Trigger apply logic - for now we can just call the command via Cobra
						_, err := ExecuteViaCobra(ctx, types.Request{Command: "apply", Args: []string{}}, os.Stdout, os.Stderr)
						if err != nil {
							slog.Error("tray apply failed", "error", err)
						}
					},
				},
				{
					Label: "Sync",
					Callback: func() {
						slog.Info("tray: triggering sync")
						_, err := ExecuteViaCobra(ctx, types.Request{Command: "sync", Args: []string{}}, os.Stdout, os.Stderr)
						if err != nil {
							slog.Error("tray sync failed", "error", err)
						}
					},
				},
				{
					Label: "Exit",
					Callback: func() {
						time.Sleep(100 * time.Millisecond) // Give time for DBus reply
						cancel()
					},
				},
			},
		})

		slog.Info("starting tray driver")
		if err := t.Run(ctx); err != nil {
			slog.Error("tray driver failed", "error", err)
		}
	}()

	listeners, err := activation.Listeners()
	if err == nil && len(listeners) > 0 {
		listener = listeners[0]
	} else {
		socketPath := getSocketPath()
		_ = os.Remove(socketPath)
		l, err := net.Listen("unix", socketPath)
		if err != nil {
			return fmt.Errorf("failed to listen on socket: %w", err)
		}
		listener = l
	}
	defer func() { _ = listener.Close() }()

	slog.Info("listening", "address", listener.Addr())

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleWS(w, r, database)
		}),
	}

	go func() {
		<-ctx.Done()
		slog.Info("context cancelled, shutting down server")
		_ = server.Close()
	}()

	return server.Serve(listener)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWS(w http.ResponseWriter, r *http.Request, database *db.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade error", "error", err)
		return
	}
	defer func() { _ = conn.Close() }()

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

	// Read loop for packets from client
	go func() {
		for {
			var packet types.StreamPacket
			if err := conn.ReadJSON(&packet); err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					slog.Debug("ws read error", "error", err)
				}
				cancel()
				return
			}

			switch packet.Type {
			case "request":
				var req types.Request
				if err := json.Unmarshal(packet.Payload, &req); err != nil {
					slog.Warn("ws unmarshal request error", "error", err)
					continue
				}
				handleRequest(ctx, req, outCh, database)
			case "history_event":
				var event types.HistoryEvent
				if err := json.Unmarshal(packet.Payload, &event); err != nil {
					slog.Warn("ws unmarshal history event error", "error", err)
					continue
				}
				if err := database.RecordHistory(ctx, event); err != nil {
					slog.Error("failed to record history", "error", err)
				}
			}
		}
	}()

	<-ctx.Done()
	<-done

	// If binary changed, exec ourselves to restart
	if shouldRestartDaemon {
		slog.Info("restarting daemon with new binary")
		exePath, err := os.Executable()
		if err != nil {
			slog.Error("failed to get executable path", "error", err)
			return
		}

		// Exec ourselves with daemon argument
		err = syscall.Exec(exePath, []string{exePath, "daemon"}, os.Environ())
		if err != nil {
			slog.Error("failed to exec daemon", "error", err)
		}
	}
}

func handleRequest(ctx context.Context, req types.Request, outCh chan types.StreamPacket, database *db.DB) {
	// Check if binary changed (mtime) - if so, signal restart needed
	if HasBinaryChanged() {
		slog.Warn("binary mtime mismatch, daemon will exec itself",
			"initial_mtime", initialMtime)

		shouldRestartDaemon = true

		resp := types.Response{
			Error: "DAEMON_RESTARTING",
		}
		payload, _ := json.Marshal(resp)
		outCh <- types.StreamPacket{
			Type:    "result",
			Payload: payload,
		}
		return
	}

	// Check if binary changed - if so, signal restart needed
	if req.BinaryHash != "" {
		daemonHash, err := exec.GetBinaryHash()
		if err == nil && daemonHash != req.BinaryHash {
			slog.Warn("binary hash mismatch, daemon will exec itself",
				"daemon_hash", daemonHash[:16],
				"client_hash", req.BinaryHash[:16])

			// Signal daemon to restart after closing this connection
			shouldRestartDaemon = true

			// Send response and close connection
			resp := types.Response{
				Error: "DAEMON_RESTARTING",
			}
			payload, _ := json.Marshal(resp)
			outCh <- types.StreamPacket{
				Type:    "result",
				Payload: payload,
			}
			return
		}
	}

	slog.Info("executing command", "command", req.Command, "args", req.Args)

	// Create logger
	handler := &logging.ChannelLogHandler{
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
	// Inject DB into context so commands can use it
	ctx = context.WithValue(ctx, types.DBKey, database)

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
	if targetCmd.DisableFlagParsing {
		argList = targetArgs
	}

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
