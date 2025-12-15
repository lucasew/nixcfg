if [[ -S "${SSH_AUTH_SOCK:-}" ]]; then
  return
fi

if [[ -S "$XDG_RUNTIME_DIR/gcr/ssh" ]]; then
  export SSH_AUTH_SOCK="$XDG_RUNTIME_DIR/gcr/ssh"
fi

