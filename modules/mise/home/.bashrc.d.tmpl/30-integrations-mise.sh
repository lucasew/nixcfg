# shellcheck shell=bash
# Mise configuration

if command -v mise >/dev/null 2>&1; then
	__ws_mise_activate="$(mise activate bash --shims 2>/dev/null)" || __ws_mise_activate=""
	if [[ -n "${__ws_mise_activate}" ]]; then
		eval "${__ws_mise_activate}"
	fi
	unset __ws_mise_activate
fi

export MISE_ALL_COMPILE=false

# Unset mise function (from mise activate) to use binary
unset -f mise 2>/dev/null || true
