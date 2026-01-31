# shellcheck shell=bash

if command -v workspaced >/dev/null 2>&1; then
	_cache_dir="${XDG_CACHE_HOME:-$HOME/.cache}/workspaced"
	mkdir -p "$_cache_dir"

	# Cache completions
	_comp_cache="$_cache_dir/completion.bash"
	if [ ! -f "$_comp_cache" ]; then
		workspaced completion bash >"$_comp_cache" 2>/dev/null
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
