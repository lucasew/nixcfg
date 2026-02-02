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
	if [ "$version" = "v2026.2.0" ]; then
		checksum_linux_x86_64="1a1455db415fa25bdfd06b76eeb9923d4b72d57efbdd8ec4d5485da6e96a4144  ./mise-v2026.2.0-linux-x64.tar.gz"
		checksum_linux_x86_64_musl="1693248fdee8dd316fc31cd64844f65fddc8778176e14842ed168de5b076d728  ./mise-v2026.2.0-linux-x64-musl.tar.gz"
		checksum_linux_arm64="ed5afe5638a6a867c9be9df2482ead1125cbdca8cb78f27a0e687ab58f4d1167  ./mise-v2026.2.0-linux-arm64.tar.gz"
		checksum_linux_arm64_musl="66236534349671d9dc9815e802b6a883c2cddb7b28d10332bc2fd992e8ed50ee  ./mise-v2026.2.0-linux-arm64-musl.tar.gz"
		checksum_linux_armv7="f3da1b8cbe387d85b8c6266a2de2a4f5d65b2bda775f9b1c9275cc1cec73a85c  ./mise-v2026.2.0-linux-armv7.tar.gz"
		checksum_linux_armv7_musl="edf0655fc232b8b4cf33818a644fa7f921fc976b1ed9a6ed18f72214953daf34  ./mise-v2026.2.0-linux-armv7-musl.tar.gz"
		checksum_macos_x86_64="0dcb83fdb17158907efdddd38acd8f1b13877f621c1d6779c8f96767daae71f1  ./mise-v2026.2.0-macos-x64.tar.gz"
		checksum_macos_arm64="f3bc94bcd49dbf4c2ac1feb0478de8a94b741d6dec8c1098c546f5e040c9eabd  ./mise-v2026.2.0-macos-arm64.tar.gz"
		checksum_linux_x86_64_zstd="8d4112a18973eea3e6559743b1b11baa68fa8ff64f432704de2cdeca742eddec  ./mise-v2026.2.0-linux-x64.tar.zst"
		checksum_linux_x86_64_musl_zstd="e26a2582ac8f92eccc156a727014cb858ce6a497ae6297b095caec19c3b443b5  ./mise-v2026.2.0-linux-x64-musl.tar.zst"
		checksum_linux_arm64_zstd="f92f7216f04138a969eeb9db3e1a751e075c2fea355acac527592cc40bbb4baa  ./mise-v2026.2.0-linux-arm64.tar.zst"
		checksum_linux_arm64_musl_zstd="e030d930235ee4d233eae70e1ea589ffca8ced21f5d9ccc626b5f8a63f20545f  ./mise-v2026.2.0-linux-arm64-musl.tar.zst"
		checksum_linux_armv7_zstd="38430652010b05b8003a5bdf7401600774ef5ed423abee68e13feb79219f2a7e  ./mise-v2026.2.0-linux-armv7.tar.zst"
		checksum_linux_armv7_musl_zstd="252fe455c03fd5146a24e35e25c46c1cf17a530d67a164fad99a684167c4ad14  ./mise-v2026.2.0-linux-armv7-musl.tar.zst"
		checksum_macos_x86_64_zstd="01d03840b3cc147cec0ed4ded1b972884e7570644878f13abc6f138cb372db47  ./mise-v2026.2.0-macos-x64.tar.zst"
		checksum_macos_arm64_zstd="4176f66ba0b30a112cd83f9b0fd06e76576a7e79e3bb21b302745bf767e6337c  ./mise-v2026.2.0-macos-arm64.tar.zst"

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
	version="${MISE_VERSION:-v2026.2.0}"
	version="${version#v}"
	os="${MISE_INSTALL_OS:-$(get_os)}"
	arch="${MISE_INSTALL_ARCH:-$(get_arch)}"
	ext="${MISE_INSTALL_EXT:-$(get_ext)}"
	install_path="${MISE_INSTALL_PATH:-$HOME/.local/bin/mise}"
	install_dir="$(dirname "$install_path")"
	install_from_github="${MISE_INSTALL_FROM_GITHUB:-}"
	if [ "$version" != "v2026.2.0" ] || [ "$install_from_github" = "1" ] || [ "$install_from_github" = "true" ]; then
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
