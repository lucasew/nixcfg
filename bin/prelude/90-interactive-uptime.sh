# shellcheck shell=bash
if [ $TERM != dumb ]; then

	if [ ! -v SD_CMD ]; then
		printf '\033[1mUptime\033[0m: %s\n' "$(uptime)" || true
	fi

fi
