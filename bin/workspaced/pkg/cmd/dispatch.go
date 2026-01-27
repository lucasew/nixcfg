package cmd

import (
	"fmt"
)

// Request definition shared for dispatch
type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type HandlerFunc func(Request) (string, error)

var handlers = map[string]HandlerFunc{
	"modn": func(req Request) (string, error) {
		return RunModn()
	},
	"media": func(req Request) (string, error) {
		return RunMedia(req.Args)
	},
	"rofi": func(req Request) (string, error) {
		return RunRofi(req.Args, req.Env)
	},
}

func ExecuteCommand(req Request) (string, error) {
	handler, ok := handlers[req.Command]
	if !ok {
		return "", fmt.Errorf("unknown command: %s", req.Command)
	}
	return handler(req)
}
