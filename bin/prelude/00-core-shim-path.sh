# shellcheck shell=bash
_shim_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)/../shim"
_workspaced_shim_dir="$HOME/.local/share/workspaced/shim/global"
_workspaced_bin_dir="$HOME/.local/share/workspaced/bin"

# Remove existing shim paths using bash string manipulation to avoid forks
_clean_path=":$PATH:"
_clean_path="${_clean_path//:$_shim_dir:/:}"
_clean_path="${_clean_path//:$_workspaced_shim_dir:/:}"
_clean_path="${_clean_path//:$_workspaced_bin_dir:/:}"
_clean_path="${_clean_path#:}"
_clean_path="${_clean_path%:}"

if [ -n "${TERMUX_VERSION:-}" ]; then
	export PATH="$_workspaced_bin_dir:$_workspaced_shim_dir:$_shim_dir:$_clean_path"
else
	export PATH="$_clean_path:$_workspaced_bin_dir:$_workspaced_shim_dir:$_shim_dir"
fi

unset _shim_dir _workspaced_shim_dir _workspaced_bin_dir _clean_path
