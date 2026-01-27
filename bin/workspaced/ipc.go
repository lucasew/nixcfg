package main

import (
	"workspaced/pkg/cmd"
)

// Request is now imported from pkg/cmd
type Request = cmd.Request

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}
