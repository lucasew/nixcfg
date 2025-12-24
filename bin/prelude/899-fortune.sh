if [ "$TERM" != dumb ]; then

	if [ ! -v SD_CMD ]; then
		timeout 1 sd fortune || true
	fi

fi
