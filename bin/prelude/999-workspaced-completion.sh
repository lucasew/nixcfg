# shellcheck shell=bash

export PATH="$HOME/.local/share/workspaced/bin:$PATH"

if command -v workspaced >/dev/null 2>&1; then
	source <(workspaced completion bash)
fi
