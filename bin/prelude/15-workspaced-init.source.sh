#!/usr/bin/env bash
# Generate workspaced init - replaces 15-workspaced-init.sh completely

if command -v workspaced >/dev/null 2>&1; then
	# Start daemon if not already running
	(workspaced daemon --try &) &>/dev/null

	# Generate completion
	workspaced completion bash

	# Apply colors directly to terminal if command exists (async to not block)
	cat <<'EOF'

	# Apply colors in interactive shell
	if [[ $- == *i* ]] && workspaced colors --help &>/dev/null; then
		(workspaced colors &) &>/dev/null
	fi
EOF
fi
