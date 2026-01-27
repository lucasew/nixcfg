package types

import "encoding/json"

type ContextKey string

const (
	DaemonModeKey ContextKey = "daemon_mode"
	EnvKey        ContextKey = "env"
	LoggerKey     ContextKey = "logger"
)

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type LogEntry struct {
	Level   string         `json:"level"`
	Message string         `json:"msg"`
	Attrs   map[string]any `json:"attrs"`
}

type StreamPacket struct {
	Type    string          `json:"type"` // "log", "result"
	Payload json.RawMessage `json:"payload"`
}
