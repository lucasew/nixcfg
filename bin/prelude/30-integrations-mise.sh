# shellcheck shell=bash
# Mise configuration - activation code generated inline in Go

eval <(mise activate || true)

export MISE_ALL_COMPILE=false

# Unset mise function (from mise activate) to use binary
unset -f mise 2>/dev/null || true
