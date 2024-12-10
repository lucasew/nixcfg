if [ $TERM != dumb ]; then

	if [ ! -v SD_CMD ]; then
		sd fortune || true
	fi

fi
