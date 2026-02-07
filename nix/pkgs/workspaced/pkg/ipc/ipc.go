package ipc

import "encoding/json"

// Request represents a generic RPC command execution request sent to the daemon.
// It carries the command to run, arguments, and any specific environment overrides.
type Request struct {
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	Env        []string `json:"env"`
	BinaryHash string   `json:"binary_hash,omitempty"` // SHA256 of client binary
}

// Response represents the final result of a command execution.
// It returns the combined output or an error message if the command failed.
type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

// StreamPacket envelopes different types of outputs to be multiplexed over a single connection.
// This allows interleaving logs, command results, and raw stdio streams.
type StreamPacket struct {
	// Type indicates the payload kind: "log", "result", "stdout", "stderr", or "history_event".
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
