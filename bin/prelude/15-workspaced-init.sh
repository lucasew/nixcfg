# shellcheck shell=bash

if command -v workspaced >/dev/null 2>&1; then
	# Start daemon if not already running
	(workspaced daemon --try &) &>/dev/null

	# Load completions
	source <(workspaced completion bash)

	# Apply colors directly to terminal if command exists
	if workspaced colors --help &>/dev/null; then
		workspaced colors
	fi
fi
