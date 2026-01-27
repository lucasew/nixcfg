package dispatch

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}
