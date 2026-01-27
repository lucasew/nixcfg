package main

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"` // Pass environment variables if needed
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}
