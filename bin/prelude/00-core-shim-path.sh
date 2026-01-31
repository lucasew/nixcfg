# shellcheck shell=bash
_shim_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)/../shim"
_workspaced_shim_dir="$HOME/.local/share/workspaced/shim/global"
_workspaced_bin_dir="$HOME/.local/share/workspaced/bin"

if [ -n "${TERMUX_VERSION:-}" ]; then
	export PATH="$_workspaced_bin_dir:$_workspaced_shim_dir:$_shim_dir:$PATH"
else
	export PATH="$PATH:$_workspaced_bin_dir:$_workspaced_shim_dir:$_shim_dir"
fi

unset _shim_dir _workspaced_shim_dir _workspaced_bin_dir _clean_path
