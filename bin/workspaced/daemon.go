package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/v22/activation"
	"workspaced/pkg/cmd"
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

	output, err := cmd.ExecuteCommand(req)
	resp := Response{Output: output}
	if err != nil {
		resp.Error = err.Error()
	}

	encoder.Encode(resp)
}
