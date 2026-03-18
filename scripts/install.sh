#!/bin/bash
set -e

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "arm" ]; then
    ARCH="arm64"
fi

REPO="DobryySoul/dockr"
BINARY="dockr"

echo "Detected system: $OS-$ARCH"

LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Could not fetch the latest release version. Make sure the repository is public and has releases."
    exit 1
fi

echo "Downloading version $LATEST_VERSION..."
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY}_${OS}_${ARCH}"

sudo curl -L "$DOWNLOAD_URL" -o /usr/local/bin/$BINARY
sudo chmod +x /usr/local/bin/$BINARY

echo "✅ Installation completed successfully!"
echo "Try running: $BINARY --help"