# shellcheck shell=bash
# Termux workaround: unset mise function that conflicts with binary

if [ -n "${TERMUX_VERSION:-}" ] && command -v mise >/dev/null; then
	unset -f mise
fi
