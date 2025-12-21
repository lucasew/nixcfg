if [ -f "$HOME/.local/bin/mise" ]; then
    MISE_CMD="$HOME/.local/bin/mise"
    if [ -n "$TERMUX_VERSION" ]; then
        MISE_CMD="termux-chroot $MISE_CMD"
    fi
    eval "$($MISE_CMD activate bash)"
fi
