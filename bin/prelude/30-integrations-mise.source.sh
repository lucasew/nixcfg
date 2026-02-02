#!/usr/bin/env bash
# Generate mise activation - replaces 30-integrations-mise.sh completely

export MISE_ALL_COMPILE=false

if [ -f "$HOME/.local/bin/mise" ]; then
	"$HOME"/.local/bin/mise activate bash

	# Add termux workaround
	cat <<'EOF'

if [ -n "${TERMUX_VERSION:-}" ] && command -v mise >/dev/null; then
	unset -f mise
fi
EOF
fi
