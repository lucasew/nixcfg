package api

import "errors"

var (
	// Driver life-cycle
	ErrDriverNotFound = errors.New("driver not found")
	ErrNotSupported   = errors.New("operation not supported")
	ErrCanceled       = errors.New("operation canceled")

	// Infrastructure
	ErrBinaryNotFound = errors.New("binary not found in PATH")
	ErrIPC            = errors.New("ipc communication failure")

	// Context/Targeting
	ErrNoFocusedOutput = errors.New("no focused output found")
	ErrNoFocusedWindow = errors.New("no focused window found")
	ErrNoTargetFound   = errors.New("target not found")
)
