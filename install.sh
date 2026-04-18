#!/bin/sh
set -eu

REPO="${REPO:-bodrovis/lokex-cli}"
BIN_NAME="${BIN_NAME:-lokex-cli}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1" >&2
    exit 1
  }
}

need_cmd uname
need_cmd mktemp
need_cmd tar
need_cmd install

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

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
  Linux)   gore_os="Linux" ;;
  Darwin)  gore_os="Darwin" ;;
  FreeBSD) gore_os="Freebsd" ;;
  *)
    echo "unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64)        gore_arch="x86_64" ;;
  i386|i686)           gore_arch="i386" ;;
  arm64|aarch64)       gore_arch="arm64" ;;
  armv6|armv6l)        gore_arch="armv6" ;;
  armv7|armv7l|armhf)  gore_arch="armv7" ;;
  *)
    echo "unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

version="${VERSION:-latest}"

if [ "$version" = "latest" ]; then
  api_url="https://api.github.com/repos/${REPO}/releases/latest"
  asset_url="$(
    fetch_stdout "$api_url" \
      | awk -F'"' '/browser_download_url/ {print $4}' \
      | grep "/${BIN_NAME}_[^/]*_${gore_os}_${gore_arch}\.tar\.gz$" \
      | head -n1 \
      || true
  )"
else
  asset_url="https://github.com/${REPO}/releases/download/${version}/${BIN_NAME}_${version}_${gore_os}_${gore_arch}.tar.gz"
fi

if [ -z "${asset_url:-}" ]; then
  echo "could not find release asset for ${gore_os}/${gore_arch}" >&2
  exit 1
fi

echo "downloading $asset_url"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

archive="$tmpdir/archive.tar.gz"

fetch_file "$asset_url" "$archive"

tar -xzf "$archive" -C "$tmpdir"

if [ ! -f "$tmpdir/$BIN_NAME" ]; then
  echo "binary not found in archive" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmpdir/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

echo
echo "installed to $INSTALL_DIR/$BIN_NAME"

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