if [ $TERM != dumb ]; then

	if [ ! -v SD_CMD ]; then
		timeout 1 figlet -f big -w $(tput cols) -c "$(hostname)" || true
	fi

fi
