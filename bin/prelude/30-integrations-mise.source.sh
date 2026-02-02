#!/usr/bin/env bash
# Generate mise activation code (output cached by workspaced shell init)

if [ -f "$HOME/.local/bin/mise" ]; then
	"$HOME"/.local/bin/mise activate bash
fi
