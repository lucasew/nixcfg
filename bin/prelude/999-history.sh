# History hook for workspaced

_workspaced_history_hook() {
	local exit_code=$?
	local cmd
	# Get last command from history
	cmd=$(HISTTIMEFORMAT= history 1 | sed 's/^[ ]*[0-9]*[ ]*//')

	# Avoid recording the record command itself and empty commands
	if [[ -z "$cmd" || "$cmd" == "workspaced dispatch history record"* ]]; then
		return
	fi

	# Send to daemon in background
	workspaced dispatch history record \
		--command "$cmd" \
		--cwd "$PWD" \
		--exit-code "$exit_code" \
		--timestamp "$(date +%s)" >/dev/null 2>&1 &
}

if [[ -n "$BASH_VERSION" ]]; then
	if [[ "$PROMPT_COMMAND" != *"_workspaced_history_hook"* ]]; then
		PROMPT_COMMAND="_workspaced_history_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
	fi
fi
