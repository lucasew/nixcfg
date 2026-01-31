# shellcheck shell=bash

if command -v workspaced >/dev/null 2>&1; then
	# Load completions directly
	source <(workspaced completion bash)

	# Apply colors directly
	eval "$(workspaced dispatch config colors)"

	# Start daemon if socket doesn't exist
	_ws_sock="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}/workspaced.sock"
	if [ ! -S "$_ws_sock" ]; then
		( workspaced daemon --try & ) &>/dev/null
	fi

	unset _ws_sock
fi

	[[ -f "$_comp_cache" ]] && source "$_comp_cache"

	# Cache colors
	_color_cache="$_cache_dir/ansi_colors.sh"
	if [ ! -f "$_color_cache" ]; then
		workspaced dispatch config colors >"$_color_cache" 2>/dev/null
	fi
	[[ -f "$_color_cache" ]] && source "$_color_cache"

	# Start daemon if socket doesn't exist
	_ws_sock="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}/workspaced.sock"
	if [ ! -S "$_ws_sock" ]; then
		(workspaced daemon --try &) &>/dev/null
	fi

	unset _cache_dir _comp_cache _color_cache _ws_sock
fi
