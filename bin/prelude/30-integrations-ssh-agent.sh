# shellcheck shell=bash
if [[ ! -S "${SSH_AUTH_SOCK:-}" ]]; then
	if [[ -S "$XDG_RUNTIME_DIR/gcr/ssh" ]]; then
		export SSH_AUTH_SOCK="$XDG_RUNTIME_DIR/gcr/ssh"
	fi
fi
