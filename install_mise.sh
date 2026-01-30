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
	if [ "$version" = "v2026.1.9" ]; then
		checksum_linux_x86_64="6577d5ac6db9b425905fa40bf27a73bf5ec0de9a45adfc59c6d3c25ed55fbd4c  ./mise-v2026.1.9-linux-x64.tar.gz"
		checksum_linux_x86_64_musl="b0b86940d6e3093a2dacce1949528651ae75168800efeef818b32f06ab29df89  ./mise-v2026.1.9-linux-x64-musl.tar.gz"
		checksum_linux_arm64="de2c8be4e77c099d1c0120465bb83ff9bce8abf10347d0766e26dd4b948d595e  ./mise-v2026.1.9-linux-arm64.tar.gz"
		checksum_linux_arm64_musl="0404770355d469757358217af9045b63fa98902c724d5da11ad76fb53a25515f  ./mise-v2026.1.9-linux-arm64-musl.tar.gz"
		checksum_linux_armv7="ec68f08417e50656e1723533022baf0b7a5d6c24c4199e46c7856a59aed1d164  ./mise-v2026.1.9-linux-armv7.tar.gz"
		checksum_linux_armv7_musl="a9532991f008eb4984c3b4c284c7ff7d8535a31dd54895c56ad614fc67b7837e  ./mise-v2026.1.9-linux-armv7-musl.tar.gz"
		checksum_macos_x86_64="972649f1d51ee1603212f180d3a9d0effd9541b23ac37c0b17d0752ede40106f  ./mise-v2026.1.9-macos-x64.tar.gz"
		checksum_macos_arm64="679a2cf4d0265582174c16f42bc947c1223a70ae57f47f2e39fe3499e443e08e  ./mise-v2026.1.9-macos-arm64.tar.gz"
		checksum_linux_x86_64_zstd="0f0f31f24f853700a4ed4656efa5175c8277c0d1d6ea64054503c4d6fe035f57  ./mise-v2026.1.9-linux-x64.tar.zst"
		checksum_linux_x86_64_musl_zstd="3bcb1aaf3f391625020a0b81ada99b4b5a5b7d0ff2d3c567d1d3d6e8c289a38d  ./mise-v2026.1.9-linux-x64-musl.tar.zst"
		checksum_linux_arm64_zstd="95e0458806c68c737b97ccdcd2b2ddbe30f874d01d37882f84b95cc4ad393ef4  ./mise-v2026.1.9-linux-arm64.tar.zst"
		checksum_linux_arm64_musl_zstd="39ff2994222f7691b9eecdab0f2858e814dd36387113f0c0cac6db2b6b69f09b  ./mise-v2026.1.9-linux-arm64-musl.tar.zst"
		checksum_linux_armv7_zstd="7300d22a01878887e650b9dcbec51615ddba5a302d65099397e1e448ea4c1695  ./mise-v2026.1.9-linux-armv7.tar.zst"
		checksum_linux_armv7_musl_zstd="7db461493b5d1ee2398b3a00fa3c2b7a55d1248586b611d89cb2e29414297944  ./mise-v2026.1.9-linux-armv7-musl.tar.zst"
		checksum_macos_x86_64_zstd="afcec62046583a6b051f32892d07456cf20ae159b7e59bafd776d4f5772ecf5b  ./mise-v2026.1.9-macos-x64.tar.zst"
		checksum_macos_arm64_zstd="172e1704b20dc52467701b6961dac3ab0cf0e42aca5b9a96b0ee32407eb97915  ./mise-v2026.1.9-macos-arm64.tar.zst"

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
	version="${MISE_VERSION:-v2026.1.9}"
	version="${version#v}"
	os="${MISE_INSTALL_OS:-$(get_os)}"
	arch="${MISE_INSTALL_ARCH:-$(get_arch)}"
	ext="${MISE_INSTALL_EXT:-$(get_ext)}"
	install_path="${MISE_INSTALL_PATH:-$HOME/.local/bin/mise}"
	install_dir="$(dirname "$install_path")"
	install_from_github="${MISE_INSTALL_FROM_GITHUB:-}"
	if [ "$version" != "v2026.1.9" ] || [ "$install_from_github" = "1" ] || [ "$install_from_github" = "true" ]; then
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
