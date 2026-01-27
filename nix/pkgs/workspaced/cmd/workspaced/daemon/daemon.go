package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/dispatch"
)

var Command = &cobra.Command{
	Use:   "daemon",
	Short: "Run the workspaced daemon",
	Run: func(c *cobra.Command, args []string) {
		if err := RunDaemon(); err != nil {
			fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
			os.Exit(1)
		}
	},
}

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
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

	fmt.Printf("Listening on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req Request
	if err := decoder.Decode(&req); err != nil {
		encoder.Encode(Response{Error: fmt.Sprintf("invalid request: %v", err)})
		return
	}

	output, err := ExecuteViaCobra(req)
	resp := Response{Output: output}
	if err != nil {
		resp.Error = err.Error()
	}

	encoder.Encode(resp)
}

var execLock sync.Mutex

func ExecuteViaCobra(req Request) (string, error) {
	execLock.Lock()
	defer execLock.Unlock()

	fullArgs := append([]string{req.Command}, req.Args...)

	root := dispatch.Command
	buf := new(bytes.Buffer)

	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(fullArgs)

	// Inject daemon_mode flag to bypass client logic
	ctx := context.WithValue(context.Background(), "env", req.Env)
	ctx = context.WithValue(ctx, "daemon_mode", true)

	err := root.ExecuteContext(ctx)
	return buf.String(), err
}
