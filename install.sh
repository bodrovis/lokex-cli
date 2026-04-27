#!/bin/sh
set -eu

REPO="${REPO:-bodrovis/lokex-cli}"
BIN_NAME="${BIN_NAME:-lokex-cli}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

need_cmd() {
	command -v "$1" >/dev/null 2>&1 || {
		echo "missing required command: $1" >&2
		exit 1
	}
}

need_cmd uname
need_cmd mktemp
need_cmd tar
need_cmd awk
need_cmd grep
need_cmd id

if command -v curl >/dev/null 2>&1; then
	fetch_stdout() {
		curl -fsSL "$1"
	}
	fetch_file() {
		curl -fsSL "$1" -o "$2"
	}
elif command -v wget >/dev/null 2>&1; then
	fetch_stdout() {
		wget -qO- "$1"
	}
	fetch_file() {
		wget -qO "$2" "$1"
	}
else
	echo "need curl or wget" >&2
	exit 1
fi

sha256_file() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | awk '{print $1}'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | awk '{print $1}'
	elif command -v openssl >/dev/null 2>&1; then
		openssl dgst -sha256 "$1" | awk '{print $NF}'
	else
		echo "need sha256sum, shasum, or openssl for checksum verification" >&2
		exit 1
	fi
}

run_as_root() {
	if [ "$(id -u)" -eq 0 ]; then
		"$@"
	elif command -v sudo >/dev/null 2>&1; then
		sudo "$@"
	else
		echo "need elevated permissions to write to $INSTALL_DIR" >&2
		echo "re-run as root or use INSTALL_DIR=\$HOME/.local/bin" >&2
		exit 1
	fi
}

install_file() {
	src="$1"
	dst="$2"

	if command -v install >/dev/null 2>&1; then
		if [ -w "$(dirname "$dst")" ]; then
			install -m 0755 "$src" "$dst"
		else
			run_as_root install -m 0755 "$src" "$dst"
		fi
	else
		need_cmd cp
		need_cmd chmod
		if [ -w "$(dirname "$dst")" ]; then
			cp "$src" "$dst"
			chmod 0755 "$dst"
		else
			run_as_root cp "$src" "$dst"
			run_as_root chmod 0755 "$dst"
		fi
	fi
}

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
Linux) gore_os="Linux" ;;
Darwin) gore_os="Darwin" ;;
FreeBSD) gore_os="Freebsd" ;;
*)
	echo "unsupported OS: $os" >&2
	exit 1
	;;
esac

case "$arch" in
x86_64 | amd64) gore_arch="x86_64" ;;
i386 | i686) gore_arch="i386" ;;
arm64 | aarch64) gore_arch="arm64" ;;
armv6 | armv6l) gore_arch="armv6" ;;
armv7 | armv7l | armhf) gore_arch="armv7" ;;
*)
	echo "unsupported architecture: $arch" >&2
	exit 1
	;;
esac

version_input="${VERSION:-latest}"

case "$version_input" in
latest)
	version="latest"
	;;
v*)
	version="$version_input"
	;;
*)
	version="v$version_input"
	;;
esac

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

if [ "$version" = "latest" ]; then
	api_url="https://api.github.com/repos/${REPO}/releases/latest"
	release_json="$tmpdir/release.json"
	fetch_file "$api_url" "$release_json"

	asset_url="$(
		awk -F'"' '/browser_download_url/ {print $4}' "$release_json" |
			grep "/${BIN_NAME}_[^/]*_${gore_os}_${gore_arch}\.tar\.gz$" |
			head -n1 ||
			true
	)"

	checksums_url="$(
		awk -F'"' '/browser_download_url/ {print $4}' "$release_json" |
			grep '/checksums\.txt$' |
			head -n1 ||
			true
	)"
else
	asset_url="https://github.com/${REPO}/releases/download/${version}/${BIN_NAME}_${version}_${gore_os}_${gore_arch}.tar.gz"
	checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
fi

if [ -z "${asset_url:-}" ]; then
	echo "could not find release asset for ${gore_os}/${gore_arch} (version=${version})" >&2
	exit 1
fi

if [ -z "${checksums_url:-}" ]; then
	echo "could not find checksums.txt for version=${version}" >&2
	exit 1
fi

asset_name="${asset_url##*/}"
archive="$tmpdir/$asset_name"
checksums_file="$tmpdir/checksums.txt"

echo "downloading $asset_name"
fetch_file "$asset_url" "$archive"
fetch_file "$checksums_url" "$checksums_file"

expected_checksum="$(
	awk -v file="$asset_name" '$2 == file { print $1; exit }' "$checksums_file"
)"

if [ -z "${expected_checksum:-}" ]; then
	echo "checksum for $asset_name not found in checksums.txt" >&2
	exit 1
fi

actual_checksum="$(sha256_file "$archive")"

if [ "$actual_checksum" != "$expected_checksum" ]; then
	echo "checksum mismatch for $asset_name" >&2
	echo "expected: $expected_checksum" >&2
	echo "actual:   $actual_checksum" >&2
	exit 1
fi

echo "checksum verified"

tar -xzf "$archive" -C "$tmpdir"

binary_path=""
if [ -f "$tmpdir/$BIN_NAME" ]; then
	binary_path="$tmpdir/$BIN_NAME"
else
	binary_path="$(find "$tmpdir" -type f -name "$BIN_NAME" | head -n1 || true)"
fi

if [ -z "$binary_path" ] || [ ! -f "$binary_path" ]; then
	echo "binary $BIN_NAME not found in archive" >&2
	exit 1
fi

if [ ! -d "$INSTALL_DIR" ]; then
	if [ -d "$(dirname "$INSTALL_DIR")" ] && [ -w "$(dirname "$INSTALL_DIR")" ]; then
		mkdir -p "$INSTALL_DIR"
	else
		run_as_root mkdir -p "$INSTALL_DIR"
	fi
fi

target="$INSTALL_DIR/$BIN_NAME"
install_file "$binary_path" "$target"

echo
echo "installed to $target"

case ":$PATH:" in
*":$INSTALL_DIR:"*)
	echo "ready to use: $BIN_NAME"
	;;
*)
	echo
	echo "add to PATH:"
	echo "export PATH=\"\$PATH:$INSTALL_DIR\""
	;;
esac
