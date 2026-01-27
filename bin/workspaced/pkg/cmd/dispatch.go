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

func ExecuteCommand(req Request) (string, error) {
	switch req.Command {
	case "modn":
		return RunModn()
	case "media":
		return RunMedia(req.Args)
	case "rofi":
		return RunRofi(req.Args, req.Env)
	default:
		return "", fmt.Errorf("unknown command: %s", req.Command)
	}
}
