# shellcheck shell=bash
mkcd() { [ -n "$1" ] && mkdir -p "$1" && cd "$_"; }
