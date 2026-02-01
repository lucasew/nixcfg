package types

import "encoding/json"

// ContextKey is a distinct string type for context values to prevent key collisions.
type ContextKey string

const (
	// DaemonModeKey flags whether the application is running as a long-lived daemon.
	DaemonModeKey ContextKey = "daemon_mode"
	// EnvKey stores the environment variables (slice of "KEY=VALUE" strings)
	// to be injected into subprocesses spawned by the handler.
	EnvKey ContextKey = "env"
	// LoggerKey stores the *slog.Logger instance, allowing contextual logging
	// (e.g., with request IDs) throughout the request lifecycle.
	LoggerKey ContextKey = "logger"
	// StdoutKey stores an io.Writer for capturing standard output from subprocesses,
	// typically redirected to the client's response stream.
	StdoutKey ContextKey = "stdout"
	// StderrKey stores an io.Writer for capturing standard error from subprocesses.
	StderrKey ContextKey = "stderr"
	// DBKey stores the *db.DB instance for database access.
	DBKey ContextKey = "db"
)

// SudoCommand defines a privileged command request.
// It aggregates execution context (cwd, env) and a timestamp for validity checks.
type SudoCommand struct {
	Slug      string   `json:"slug"`
	Command   string   `json:"command"`
	Args      []string `json:"args"`
	Cwd       string   `json:"cwd"`
	Env       []string `json:"env"`
	Timestamp int64    `json:"timestamp"`
}

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

// LogEntry is a serializable representation of a structured log record.
// It is used to marshal log events to JSON for transmission over the wire.
type LogEntry struct {
	Level   string         `json:"level"`
	Message string         `json:"msg"`
	Attrs   map[string]any `json:"attrs"`
}

// StreamPacket envelopes different types of outputs to be multiplexed over a single connection.
// This allows interleaving logs, command results, and raw stdio streams.
type StreamPacket struct {
	// Type indicates the payload kind: "log", "result", "stdout", "stderr", or "history_event".
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// HistoryEvent is sent by the shell hook to record a command execution.
type HistoryEvent struct {
	Command   string `json:"command"`
	Cwd       string `json:"cwd"`
	Timestamp int64  `json:"timestamp"`
	ExitCode  int    `json:"exit_code"`
	Duration  int64  `json:"duration_ms"`
}
