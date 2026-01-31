# shellcheck shell=bash

if command -v workspaced >/dev/null 2>&1; then
	# Load completions
	source <(workspaced completion bash)

	# Apply colors
	source <(workspaced dispatch config colors)

	# Start daemon if socket doesn't exist
	_ws_sock="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}/workspaced.sock"
	if [ ! -S "$_ws_sock" ]; then
		(workspaced daemon --try &) &>/dev/null
	fi

	unset _ws_sock
fi
