if [ ! -v SD_CMD ]; then
figlet -f big -w $(tput cols) -c "$(hostname)" || true
fi

