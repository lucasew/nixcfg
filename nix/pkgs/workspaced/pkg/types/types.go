package types

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type ContextKey string

const (
	DaemonModeKey ContextKey = "daemon_mode"
	EnvKey        ContextKey = "env"
)
