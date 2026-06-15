#!/bin/bash
set -euo pipefail

BUFFMAN_VERSION="${BUFFMAN_VERSION:-1.0.0}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS=$(uname -s)
case "$OS" in
    Linux)  OS_NAME="linux" ;;
    Darwin) OS_NAME="darwin" ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  ARCH_NAME="amd64" ;;
    arm64 | aarch64) ARCH_NAME="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

ARCHIVE="buffman_${BUFFMAN_VERSION}_${OS_NAME}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/the-protobuf-project/buffman/releases/download/v${BUFFMAN_VERSION}/${ARCHIVE}"
TMP_DIR=$(mktemp -d)

echo "Downloading buffman v${BUFFMAN_VERSION} for ${OS_NAME}/${ARCH_NAME}..."

if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}"; then
    echo "Failed to download buffman from: $DOWNLOAD_URL"
    rm -rf "$TMP_DIR"
    exit 1
fi

tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"
sudo install -m 0755 "${TMP_DIR}/buffman" "${INSTALL_DIR}/buffman"

rm -rf "$TMP_DIR"

echo "buffman v${BUFFMAN_VERSION} installed to ${INSTALL_DIR}/buffman"
