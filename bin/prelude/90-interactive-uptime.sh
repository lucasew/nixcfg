# shellcheck shell=bash
if [[ $- == *i* ]]; then
	printf '\033[1mUptime\033[0m: %s\n' "$(uptime)" || true
fi
