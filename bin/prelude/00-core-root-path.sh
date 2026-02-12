# shellcheck shell=bash
# Get the root path without subshells
if [ -z "${NIXCFG_ROOT_PATH:-}" ]; then
	_script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
	# If we are in process substitution, this might fail or be wrong
	if [[ "$_script_dir" == "/dev/fd/"* ]] || [[ "$_script_dir" == "." ]]; then
		# Fallback to current directory or common locations if really lost
		# But usually NIXCFG_ROOT_PATH is injected by workspaced generator
		export NIXCFG_ROOT_PATH="$HOME/.dotfiles"
	else
		export NIXCFG_ROOT_PATH="$(realpath "$_script_dir/../..")"
	fi
	unset _script_dir
fi
