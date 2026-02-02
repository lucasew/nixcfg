# shellcheck shell=bash
# Termux workaround: unset mise function that conflicts with binary

if command -v mise >/dev/null; then
	unset -f mise
fi
