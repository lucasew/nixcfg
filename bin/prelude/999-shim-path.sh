# shellcheck shell=bash
_shim_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)/../shim"
_workspaced_shim_dir="$HOME/.local/share/workspaced/shim/global"
_path_without_shim="$(echo "$PATH" | tr ':' '\n' | grep -v 'bin/shim' | grep -v '.local/share/workspaced/shim/global' | grep . | paste -sd : -)"

if [ -n "${TERMUX_VERSION:-}" ]; then
	export PATH="$_workspaced_shim_dir:$_shim_dir:$_path_without_shim"
else
	export PATH="$_path_without_shim:$_workspaced_shim_dir:$_shim_dir"
fi

unset _shim_dir _workspaced_shim_dir _path_without_shim
