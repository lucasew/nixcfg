# shellcheck shell=bash
# Get the root path without subshells
if [ -z "${NIXCFG_ROOT_PATH:-}" ]; then
	_script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
	export NIXCFG_ROOT_PATH="$(realpath "$_script_dir/../..")"
	unset _script_dir
fi
