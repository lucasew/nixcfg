# shellcheck shell=bash
# This file is replaced by 30-integrations-mise.source.sh output when using workspaced shell init
# The .source.sh file executes mise activate bash in parallel with other commands for faster startup

export MISE_ALL_COMPILE=false

# Fallback for when workspaced is not available
if [ ! -v WORKSPACED_SHELL_INIT ]; then
	if [ -f "$HOME/.local/bin/mise" ]; then
		eval "$("$HOME"/.local/bin/mise activate bash)"
	fi
fi

if [ -n "${TERMUX_VERSION:-}" ] && command -v mise >/dev/null; then
	unset -f mise
fi
