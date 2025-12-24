if [ -f "$HOME/.local/bin/mise" ]; then
    if command -v mise >/dev/null 2>&1; then
        eval "$(mise activate bash)"
    else
        eval "$("$HOME"/.local/bin/mise activate bash)"
    fi
fi
