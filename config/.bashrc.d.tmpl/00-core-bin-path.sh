# shellcheck shell=bash
_local_bin="$HOME/.local/bin"

# Remove existing ~/.local/bin from PATH to avoid duplicates
_clean_path=":$PATH:"
_clean_path="${_clean_path//:$_local_bin:/:}"
_clean_path="${_clean_path#:}"
_clean_path="${_clean_path%:}"

# Prepend ~/.local/bin (global shims including 'x' helper)
export PATH="$_local_bin:$_clean_path"

unset _local_bin _clean_path

# Note: Lazy shims in ~/.local/share/workspaced/shim/lazy/
# are added to PATH via 'x' wrapper when needed
