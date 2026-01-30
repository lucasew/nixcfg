# shellcheck shell=bash
alias la='ls -lha'
alias l='ls'
alias cd..='cd ..'
alias ..='cd ..'
alias ç='sd'
alias e=$EDITOR
alias sdw=sd
alias workspaced="$HOME/.local/share/workspaced/bin/workspaced"

function reset_term {
	tput reset
	type -t setup_colors >/dev/null && setup_colors
	if [ -v PREFIX ]; then
		cat $PREFIX/etc/motd
	else
		for item in /etc/motd*; do
			if [ ! -d "$(realpath "$item")" ]; then
				cat "$item"
			fi
		done
	fi
}
