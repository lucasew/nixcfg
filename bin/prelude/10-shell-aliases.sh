# shellcheck shell=bash
alias la='ls -lha'
alias l='ls'
alias cd..='cd ..'
alias ..='cd ..'
alias ç='sd'
export EDITOR=${EDITOR:-hx}
alias e=$EDITOR
alias sdw=sd

function reset_term {
	tput reset
	workspaced colors
}
