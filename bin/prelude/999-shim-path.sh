# shellcheck shell=bash
_shim_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)/../shim"
_path_without_shim="$(echo "$PATH" | tr ':' '\n' | grep -v 'bin/shim' | grep . | paste -sd : -)"

if [ -n "${TERMUX_VERSION:-}" ]; then
	export PATH="$_shim_dir:$_path_without_shim"
else
	export PATH="$_path_without_shim:$_shim_dir"
fi

unset _shim_dir _path_without_shim
