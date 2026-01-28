#!/bin/sh
set -eu

#region logging setup
if [ "${MISE_DEBUG-}" = "true" ] || [ "${MISE_DEBUG-}" = "1" ]; then
	debug() {
		echo "$@" >&2
	}
else
	debug() {
		:
	}
fi

if [ "${MISE_QUIET-}" = "1" ] || [ "${MISE_QUIET-}" = "true" ]; then
	info() {
		:
	}
else
	info() {
		echo "$@" >&2
	}
fi

error() {
	echo "$@" >&2
	exit 1
}
#endregion

#region environment setup
get_os() {
	os="$(uname -s)"
	if [ "$os" = Darwin ]; then
		echo "macos"
	elif [ "$os" = Linux ]; then
		echo "linux"
	else
		error "unsupported OS: $os"
	fi
}

get_arch() {
	musl=""
	if type ldd >/dev/null 2>/dev/null; then
		if [ "${MISE_INSTALL_MUSL-}" = "1" ] || [ "${MISE_INSTALL_MUSL-}" = "true" ]; then
			musl="-musl"
		elif [ "$(uname -o)" = "Android" ]; then
			# Android (Termux) always uses musl
			musl="-musl"
		else
			libc=$(ldd /bin/ls | grep 'musl' | head -1 | cut -d ' ' -f1)
			if [ -n "$libc" ]; then
				musl="-musl"
			fi
		fi
	fi
	arch="$(uname -m)"
	if [ "$arch" = x86_64 ]; then
		echo "x64$musl"
	elif [ "$arch" = aarch64 ] || [ "$arch" = arm64 ]; then
		echo "arm64$musl"
	elif [ "$arch" = armv7l ]; then
		echo "armv7$musl"
	else
		error "unsupported architecture: $arch"
	fi
}

get_ext() {
	if [ -n "${MISE_INSTALL_EXT:-}" ]; then
		echo "$MISE_INSTALL_EXT"
	elif [ -n "${MISE_VERSION:-}" ] && echo "$MISE_VERSION" | grep -q '^v2024'; then
		# 2024 versions don't have zstd tarballs
		echo "tar.gz"
	elif tar_supports_zstd; then
		echo "tar.zst"
	elif command -v zstd >/dev/null 2>&1; then
		echo "tar.zst"
	else
		echo "tar.gz"
	fi
}

tar_supports_zstd() {
	# tar is bsdtar or version is >= 1.31
	if tar --version | grep -q 'bsdtar' && command -v zstd >/dev/null 2>&1; then
		true
	elif tar --version | grep -q '1\.(3[1-9]|[4-9][0-9]'; then
		true
	else
		false
	fi
}

shasum_bin() {
	if command -v shasum >/dev/null 2>&1; then
		echo "shasum"
	elif command -v sha256sum >/dev/null 2>&1; then
		echo "sha256sum"
	else
		error "mise install requires shasum or sha256sum but neither is installed. Aborting."
	fi
}

get_checksum() {
	version=$1
	os=$2
	arch=$3
	ext=$4
	url="https://github.com/jdx/mise/releases/download/v${version}/SHASUMS256.txt"

	# For current version use static checksum otherwise
	# use checksum from releases
	if [ "$version" = "v2026.1.8" ]; then
		checksum_linux_x86_64="bc19876755b7cc91801009d98dad742105829ca546c45a815478d627b37a4a25  ./mise-v2026.1.8-linux-x64.tar.gz"
		checksum_linux_x86_64_musl="1ef3aaefd4aabf186f85d0aacf2da77d8acf8589763b0d573cb30a8bf094ede1  ./mise-v2026.1.8-linux-x64-musl.tar.gz"
		checksum_linux_arm64="1b000dd2a3b698ec8bf884130f923ae89de91f3ce8bdc0422bbb0f2a6b539061  ./mise-v2026.1.8-linux-arm64.tar.gz"
		checksum_linux_arm64_musl="9df08cf571d61c70e12601a6bb8f3663559672cd43413ac5937b4a056d0701d8  ./mise-v2026.1.8-linux-arm64-musl.tar.gz"
		checksum_linux_armv7="187eb0653a2c90276dd38a2a62d8c85c3248fd7d9cbdd1bb26b7e2fc52d6d20c  ./mise-v2026.1.8-linux-armv7.tar.gz"
		checksum_linux_armv7_musl="1ac7a9253c121aed1aaba7a1d581e6bb76addc2d072e0b6f1ac041fa81204794  ./mise-v2026.1.8-linux-armv7-musl.tar.gz"
		checksum_macos_x86_64="f72b5be1f81709c22b0b3bda9eecd58c9e368af64a385d58dcd8a564e46f8cfb  ./mise-v2026.1.8-macos-x64.tar.gz"
		checksum_macos_arm64="8f858cbc78131850f507d5fea49c892f49e9cc3793f45cd7b7feea37fd00ec7c  ./mise-v2026.1.8-macos-arm64.tar.gz"
		checksum_linux_x86_64_zstd="0614738c61260df55cc3435a7e1bcd88bd627229cefd39c5da4318cd7b0b09af  ./mise-v2026.1.8-linux-x64.tar.zst"
		checksum_linux_x86_64_musl_zstd="68f8d9c1d08d5412a5281fde71865f267d9a8fd9b755eacd5a33959e5a71a3d2  ./mise-v2026.1.8-linux-x64-musl.tar.zst"
		checksum_linux_arm64_zstd="aa5818d14dc27b00d29b7a1981b5a2bf7295dd44b0672061abe8444a44a0c0f3  ./mise-v2026.1.8-linux-arm64.tar.zst"
		checksum_linux_arm64_musl_zstd="a880f7ede396ad67e4e37b88520adc10ff3ef6643c877aabda04b6f9bfbf741b  ./mise-v2026.1.8-linux-arm64-musl.tar.zst"
		checksum_linux_armv7_zstd="21b38392a73905f456318904687b7c9d2ce5611405780d04b493e49f6e264485  ./mise-v2026.1.8-linux-armv7.tar.zst"
		checksum_linux_armv7_musl_zstd="28bd2a39ba35d2a371e9fcaa0ee625726311916dffaac6f2b956a1efed137b4f  ./mise-v2026.1.8-linux-armv7-musl.tar.zst"
		checksum_macos_x86_64_zstd="48aa1a63e3917d2dd74028fc5aacf4b5bac8fa1ddf203c821d699522e1ac3e3f  ./mise-v2026.1.8-macos-x64.tar.zst"
		checksum_macos_arm64_zstd="0f4fec2a6711d37b360041a181e2370279740d02c13d363d19b1362686cbd3a9  ./mise-v2026.1.8-macos-arm64.tar.zst"

		# TODO: refactor this, it's a bit messy
		if [ "$ext" = "tar.zst" ]; then
			if [ "$os" = "linux" ]; then
				if [ "$arch" = "x64" ]; then
					echo "$checksum_linux_x86_64_zstd"
				elif [ "$arch" = "x64-musl" ]; then
					echo "$checksum_linux_x86_64_musl_zstd"
				elif [ "$arch" = "arm64" ]; then
					echo "$checksum_linux_arm64_zstd"
				elif [ "$arch" = "arm64-musl" ]; then
					echo "$checksum_linux_arm64_musl_zstd"
				elif [ "$arch" = "armv7" ]; then
					echo "$checksum_linux_armv7_zstd"
				elif [ "$arch" = "armv7-musl" ]; then
					echo "$checksum_linux_armv7_musl_zstd"
				else
					warn "no checksum for $os-$arch"
				fi
			elif [ "$os" = "macos" ]; then
				if [ "$arch" = "x64" ]; then
					echo "$checksum_macos_x86_64_zstd"
				elif [ "$arch" = "arm64" ]; then
					echo "$checksum_macos_arm64_zstd"
				else
					warn "no checksum for $os-$arch"
				fi
			else
				warn "no checksum for $os-$arch"
			fi
		else
			if [ "$os" = "linux" ]; then
				if [ "$arch" = "x64" ]; then
					echo "$checksum_linux_x86_64"
				elif [ "$arch" = "x64-musl" ]; then
					echo "$checksum_linux_x86_64_musl"
				elif [ "$arch" = "arm64" ]; then
					echo "$checksum_linux_arm64"
				elif [ "$arch" = "arm64-musl" ]; then
					echo "$checksum_linux_arm64_musl"
				elif [ "$arch" = "armv7" ]; then
					echo "$checksum_linux_armv7"
				elif [ "$arch" = "armv7-musl" ]; then
					echo "$checksum_linux_armv7_musl"
				else
					warn "no checksum for $os-$arch"
				fi
			elif [ "$os" = "macos" ]; then
				if [ "$arch" = "x64" ]; then
					echo "$checksum_macos_x86_64"
				elif [ "$arch" = "arm64" ]; then
					echo "$checksum_macos_arm64"
				else
					warn "no checksum for $os-$arch"
				fi
			else
				warn "no checksum for $os-$arch"
			fi
		fi
	else
		if command -v curl >/dev/null 2>&1; then
			debug ">" curl -fsSL "$url"
			checksums="$(curl --compressed -fsSL "$url")"
		else
			if command -v wget >/dev/null 2>&1; then
				debug ">" wget -qO - "$url"
				checksums="$(wget -qO - "$url")"
			else
				error "mise standalone install specific version requires curl or wget but neither is installed. Aborting."
			fi
		fi
		# TODO: verify with minisign or gpg if available

		checksum="$(echo "$checksums" | grep "$os-$arch.$ext")"
		if ! echo "$checksum" | grep -Eq "^([0-9a-f]{32}|[0-9a-f]{64})"; then
			warn "no checksum for mise $version and $os-$arch"
		else
			echo "$checksum"
		fi
	fi
}

#endregion

download_file() {
	url="$1"
	download_dir="$2"
	filename="$(basename "$url")"
	file="$download_dir/$filename"

	info "mise: installing mise..."

	if command -v curl >/dev/null 2>&1; then
		debug ">" curl -#fLo "$file" "$url"
		curl -#fLo "$file" "$url"
	else
		if command -v wget >/dev/null 2>&1; then
			debug ">" wget -qO "$file" "$url"
			stderr=$(mktemp)
			wget -O "$file" "$url" >"$stderr" 2>&1 || error "wget failed: $(cat "$stderr")"
			rm "$stderr"
		else
			error "mise standalone install requires curl or wget but neither is installed. Aborting."
		fi
	fi

	echo "$file"
}

install_mise() {
	version="${MISE_VERSION:-v2026.1.8}"
	version="${version#v}"
	os="${MISE_INSTALL_OS:-$(get_os)}"
	arch="${MISE_INSTALL_ARCH:-$(get_arch)}"
	ext="${MISE_INSTALL_EXT:-$(get_ext)}"
	install_path="${MISE_INSTALL_PATH:-$HOME/.local/bin/mise}"
	install_dir="$(dirname "$install_path")"
	install_from_github="${MISE_INSTALL_FROM_GITHUB:-}"
	if [ "$version" != "v2026.1.8" ] || [ "$install_from_github" = "1" ] || [ "$install_from_github" = "true" ]; then
		tarball_url="https://github.com/jdx/mise/releases/download/v${version}/mise-v${version}-${os}-${arch}.${ext}"
	elif [ -n "${MISE_TARBALL_URL-}" ]; then
		tarball_url="$MISE_TARBALL_URL"
	else
		tarball_url="https://mise.jdx.dev/v${version}/mise-v${version}-${os}-${arch}.${ext}"
	fi

	download_dir="$(mktemp -d)"
	cache_file=$(download_file "$tarball_url" "$download_dir")
	debug "mise-setup: tarball=$cache_file"

	debug "validating checksum"
	cd "$(dirname "$cache_file")" && get_checksum "$version" "$os" "$arch" "$ext" | "$(shasum_bin)" -c >/dev/null

	# extract tarball
	mkdir -p "$install_dir"
	rm -rf "$install_path"
	extract_dir="$(mktemp -d)"
	cd "$extract_dir"
	if [ "$ext" = "tar.zst" ] && ! tar_supports_zstd; then
		zstd -d -c "$cache_file" | tar -xf -
	else
		tar -xf "$cache_file"
	fi
	mv mise/bin/mise "$install_path"

	# cleanup
	cd / # Move out of $extract_dir before removing it
	rm -rf "$download_dir"
	rm -rf "$extract_dir"

	info "mise: installed successfully to $install_path"
}

after_finish_help() {
	case "${SHELL:-}" in
	*/zsh)
		info "mise: run the following to activate mise in your shell:"
		info "echo \"eval \\\"\\\$($install_path activate zsh)\\\"\" >> \"${ZDOTDIR-$HOME}/.zshrc\""
		info ""
		info "mise: run \`mise doctor\` to verify this is setup correctly"
		;;
	*/bash)
		info "mise: run the following to activate mise in your shell:"
		info "echo \"eval \\\"\\\$($install_path activate bash)\\\"\" >> ~/.bashrc"
		info ""
		info "mise: run \`mise doctor\` to verify this is setup correctly"
		;;
	*/fish)
		info "mise: run the following to activate mise in your shell:"
		info "echo \"$install_path activate fish | source\" >> ~/.config/fish/config.fish"
		info ""
		info "mise: run \`mise doctor\` to verify this is setup correctly"
		;;
	*)
		info "mise: run \`$install_path --help\` to get started"
		;;
	esac
}

install_mise
if [ "${MISE_INSTALL_HELP-}" != 0 ]; then
	after_finish_help
fi
