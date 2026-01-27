package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// TryRemote attempts to run the command via the daemon.
// Returns (output, true, nil) if successful.
// Returns ("", false, nil) if daemon is unreachable.
// Returns ("", true, error) if daemon was reached but returned an error.
func TryRemote(cmd string, args []string) (string, bool, error) {
	socketPath := getSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return "", false, nil // Daemon likely not running
	}
	defer conn.Close()

	req := Request{
		Command: cmd,
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
