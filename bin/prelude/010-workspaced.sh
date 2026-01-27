function workspaced {
	# Check if there is a real workspaced binary in the PATH that is not this function
	local real_workspaced
	real_workspaced=$(type -ap workspaced | grep -v "function" | head -n 1)

	if [[ -n "$real_workspaced" ]]; then
		"$real_workspaced" "$@"
	else
		# Fallback to hot build and run from dotfiles
		local dotfiles_root
		dotfiles_root=$(sd d root)
		local source_dir="$dotfiles_root/nix/pkgs/workspaced"

		if [[ -d "$source_dir" ]]; then
			local temp_bin="${XDG_RUNTIME_DIR:-/tmp}/workspaced-hot"
			(cd "$source_dir" && go build -o "$temp_bin" ./cmd/workspaced)
			"$temp_bin" "$@"
		else
			echo "workspaced: binary not found and source not found in $dotfiles_root" >&2
			return 1
		fi
	fi
}
