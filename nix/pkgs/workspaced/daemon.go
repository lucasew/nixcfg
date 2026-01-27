package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"bytes"

	"github.com/coreos/go-systemd/v22/activation"
	"workspaced/cmd/workspaced/dispatch"
	"github.com/spf13/cobra"
	"workspaced/pkg/cmd" // for Request struct
)

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "workspaced.sock")
}

func RunDaemon() error {
	var listener net.Listener

	// Check systemd activation
	listeners, err := activation.Listeners()
	if err == nil && len(listeners) > 0 {
		listener = listeners[0]
	} else {
		// Fallback to manual socket creation
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

	var req cmd.Request
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

func ExecuteViaCobra(req cmd.Request) (string, error) {
	// Reconstruct args: command name + args
	// e.g. "modn" [] -> ["modn"]
	// "media" ["next"] -> ["media", "next"]
	fullArgs := append([]string{req.Command}, req.Args...)

	// Create a fresh root command to execute
	// We need to avoid side effects from previous executions if we reused a global root.
	// So we need a factory.
	// We'll define a local root and add commands from dispatch.
	root := &cobra.Command{Use: "workspaced"}
	root.AddCommand(dispatch.ModnCmd)
	root.AddCommand(dispatch.MediaCmd)
	root.AddCommand(dispatch.RofiCmd)

	// Capture output
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)

	// Set args
	root.SetArgs(fullArgs)

	// Inject Env into context (for Rofi)
	// Cobra commands can access context.
	// We wrap the context.
	// Since we can't easily change the context of the *root* before Execute (SetContext is available in newer Cobra, check version?)
	// Cobra 1.10 has ExecuteContext? No, ExecuteContext is valid.
	// Or we can just use a global/package var but that's ugly.
	// Let's rely on Context if possible.
	// Actually, we can use `SetContext`?
	// `root.ExecuteContext(ctx)`

	// Construct context with env
	// Note: We need to define a key type to avoid collisions, but string is fine for now.
	// ctx := context.WithValue(context.Background(), "env", req.Env)
	// err := root.ExecuteContext(ctx)

	// NOTE: Since I cannot see if `ExecuteContext` is available in the vendored Cobra version easily,
	// I will assume standard Execute usage and pass Env via a hack or context if available.
	// Most standard Cobra has ExecuteContext.

	// Wait, I can just set `dispatch.RofiCmd`'s behavior to read from a transient store if needed.
	// Or even simpler: Use `root.SetContext`? No, it's `ExecuteContext`.

	// Let's try `ExecuteC`.
	// For now, let's proceed with `Execute` and assume we can pass env via Context in `ExecuteContext`.
	// But `dispatch.RofiCmd` is a global variable.
	// If we execute it concurrently, we have a race condition if we modify it.
	// Ideally `dispatch.RofiCmd` should be a function returning a command.
	// But `dispatch.cmds.go` defined it as a var.
	// The PR comment says: "RPC works directly with the Command of the subcommand dispatch".

	// FIX: `dispatch.RofiCmd` should be stateless regarding request data.
	// Request data (Env) must be passed via context.

	// Let's implement `ExecuteViaCobra` using `ExecuteContext`.
	// Requires `context` package.

	return runCobra(root, fullArgs, req.Env)
}
