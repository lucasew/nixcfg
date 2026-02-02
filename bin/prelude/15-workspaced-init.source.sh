#!/usr/bin/env bash
# Generate workspaced init - replaces 15-workspaced-init.sh completely

if command -v workspaced >/dev/null 2>&1; then
	# Start daemon if not already running
	(workspaced daemon --try &) &>/dev/null

	# Generate completion
	workspaced completion bash

	# Apply colors in interactive shell
	echo ""
	echo "if [[ \$- == *i* ]] && workspaced colors --help &>/dev/null; then"
	echo "	workspaced colors"
	echo "fi"
fi
