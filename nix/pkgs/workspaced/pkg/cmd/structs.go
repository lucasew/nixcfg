package cmd

// Request definition shared for dispatch
type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}
