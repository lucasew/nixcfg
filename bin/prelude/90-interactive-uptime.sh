# shellcheck shell=bash
if [ $TERM != dumb ]; then
	printf '\033[1mUptime\033[0m: %s\n' "$(uptime)" || true
fi
