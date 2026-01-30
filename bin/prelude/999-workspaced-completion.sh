# shellcheck shell=bash

if command -v workspaced >/dev/null 2>&1; then
	source <(workspaced completion bash)
fi
