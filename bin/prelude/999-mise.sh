export MISE_ALL_COMPILE=false

if [ -f "$HOME/.local/bin/mise" ]; then
    eval "$("$HOME"/.local/bin/mise activate bash)"
    
fi
