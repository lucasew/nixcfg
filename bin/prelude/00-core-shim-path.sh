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

# Always prepend shims to ensure they can manage binary updates
export PATH="$_shim_dir:$_workspaced_shim_dir:$_workspaced_bin_dir:$_clean_path"

unset _shim_dir _workspaced_shim_dir _workspaced_bin_dir _clean_path
