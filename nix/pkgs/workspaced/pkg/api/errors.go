package api

import "errors"

var (
	// Driver life-cycle
	ErrDriverNotFound = errors.New("screen driver not found: WAYLAND_DISPLAY or DISPLAY environment variable not set. Try: export WAYLAND_DISPLAY=wayland-1")
	ErrNotSupported   = errors.New("operation not supported")
	ErrCanceled       = errors.New("operation canceled")
	ErrNotImplemented = errors.New("not implemented")

	// Infrastructure
	ErrBinaryNotFound = errors.New("binary not found in PATH")
	ErrIPC            = errors.New("ipc communication failure")
	ErrNetwork        = errors.New("network failure")

	// Context/Targeting
	ErrNoFocusedOutput = errors.New("no focused output found")
	ErrNoFocusedWindow = errors.New("no focused window found")
	ErrNoTargetFound   = errors.New("target not found")
	ErrNoTerminalFound = errors.New("no terminal found")

	// Configuration
	ErrConfigNotFound       = errors.New("config not found")
	ErrHostNotFound         = errors.New("host not found")
	ErrInvalidAddr          = errors.New("invalid address")
	ErrDotfilesRootNotFound = errors.New("dotfiles root not found")

	// Nix/Build
	ErrBuildFailed      = errors.New("build failed")
	ErrFlakeRefRequired = errors.New("flake reference required")
)
