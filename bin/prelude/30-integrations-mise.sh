# shellcheck shell=bash
export MISE_ALL_COMPILE=false

if [ -f "$HOME/.local/bin/mise" ]; then
	eval "$("$HOME"/.local/bin/mise activate bash)"
	if [ -n "${TERMUX_VERSION:-}" ]; then
		unset -f mise
	fi
fi
