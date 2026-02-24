# shellcheck shell=bash
# Mise configuration

_mise_shims_dir="$HOME/.local/share/mise/shims"
_clean_path=":$PATH:"
_clean_path="${_clean_path//:$_mise_shims_dir:/:}"
_clean_path="${_clean_path#:}"
_clean_path="${_clean_path%:}"
export PATH="${_clean_path}"
unset _mise_shims_dir _clean_path

if command -v mise >/dev/null 2>&1; then
	__ws_mise_activate="$(mise activate bash 2>/dev/null)" || __ws_mise_activate=""
	if [[ -n "${__ws_mise_activate}" ]]; then
		eval "${__ws_mise_activate}" 2>/dev/null
	fi
	unset __ws_mise_activate
fi

export MISE_ALL_COMPILE=false
