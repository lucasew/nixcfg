# shellcheck shell=bash

# Apply base16 color scheme from settings.toml

function setup_colors {
	if ! command -v workspaced >/dev/null 2>&1; then
		return
	fi

	# Read base16 colors from settings.toml
	local palette
	palette=$(workspaced dispatch config get desktop.palette.base16 2>/dev/null) || return

	# Extract colors using jq if available, otherwise skip
	if ! command -v jq >/dev/null 2>&1; then
		return
	fi

	# Map base16 to ANSI colors
	local base00 base01 base02 base03 base04 base05 base06 base07
	local base08 base09 base0A base0B base0C base0D base0E base0F

	base00=$(echo "$palette" | jq -r '.base00 // empty')
	base01=$(echo "$palette" | jq -r '.base01 // empty')
	base02=$(echo "$palette" | jq -r '.base02 // empty')
	base03=$(echo "$palette" | jq -r '.base03 // empty')
	base04=$(echo "$palette" | jq -r '.base04 // empty')
	base05=$(echo "$palette" | jq -r '.base05 // empty')
	base06=$(echo "$palette" | jq -r '.base06 // empty')
	base07=$(echo "$palette" | jq -r '.base07 // empty')
	base08=$(echo "$palette" | jq -r '.base08 // empty')
	base09=$(echo "$palette" | jq -r '.base09 // empty')
	base0A=$(echo "$palette" | jq -r '.base0A // empty')
	base0B=$(echo "$palette" | jq -r '.base0B // empty')
	base0C=$(echo "$palette" | jq -r '.base0C // empty')
	base0D=$(echo "$palette" | jq -r '.base0D // empty')
	base0E=$(echo "$palette" | jq -r '.base0E // empty')
	base0F=$(echo "$palette" | jq -r '.base0F // empty')

	[[ -z "$base00" ]] && return

	# Apply colors using ANSI escape sequences
	# Standard ANSI color mapping
	printf '\033]4;0;#%s\033\\' "$base00"  # black
	printf '\033]4;1;#%s\033\\' "$base08"  # red
	printf '\033]4;2;#%s\033\\' "$base0B"  # green
	printf '\033]4;3;#%s\033\\' "$base0A"  # yellow
	printf '\033]4;4;#%s\033\\' "$base0D"  # blue
	printf '\033]4;5;#%s\033\\' "$base0E"  # magenta
	printf '\033]4;6;#%s\033\\' "$base0C"  # cyan
	printf '\033]4;7;#%s\033\\' "$base05"  # white
	printf '\033]4;8;#%s\033\\' "$base03"  # bright black
	printf '\033]4;9;#%s\033\\' "$base08"  # bright red
	printf '\033]4;10;#%s\033\\' "$base0B" # bright green
	printf '\033]4;11;#%s\033\\' "$base0A" # bright yellow
	printf '\033]4;12;#%s\033\\' "$base0D" # bright blue
	printf '\033]4;13;#%s\033\\' "$base0E" # bright magenta
	printf '\033]4;14;#%s\033\\' "$base0C" # bright cyan
	printf '\033]4;15;#%s\033\\' "$base07" # bright white

	# Set foreground and background
	printf '\033]10;#%s\033\\' "$base05"   # foreground
	printf '\033]11;#%s\033\\' "$base00"   # background
	printf '\033]12;#%s\033\\' "$base05"   # cursor
}

setup_colors
