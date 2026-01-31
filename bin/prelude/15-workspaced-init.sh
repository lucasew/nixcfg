# shellcheck shell=bash

if command -v workspaced >/dev/null 2>&1; then
	(workspaced daemon --try &) &> /dev/null
	# Load completions and colors in one go
	source <(
		workspaced completion bash;
		workspaced dispatch config colors
	)
fi
