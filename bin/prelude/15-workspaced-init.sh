# shellcheck shell=bash
# This file is partially replaced by 15-workspaced-init.source.sh when using workspaced shell init
# The .source.sh generates completion code in parallel for faster startup

if command -v workspaced >/dev/null 2>&1; then
	# Start daemon if not already running
	(workspaced daemon --try &) &>/dev/null

	# Load completions - fallback when not using workspaced shell init
	if [ ! -v WORKSPACED_SHELL_INIT ]; then
		source <(workspaced completion bash)
	fi

	# Apply colors directly to terminal if command exists (async to not block)
	if [[ $- == *i* ]] && workspaced colors --help &>/dev/null; then
		(workspaced colors &) &>/dev/null
	fi
fi
