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

echo "Определена система: $OS-$ARCH"

LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Не удалось получить информацию о последней версии. Убедитесь, что репозиторий публичный и имеет релизы."
    exit 1
fi

echo "Скачивание версии $LATEST_VERSION..."
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY}_${OS}_${ARCH}"

sudo curl -L "$DOWNLOAD_URL" -o /usr/local/bin/$BINARY
sudo chmod +x /usr/local/bin/$BINARY

echo "✅ Установка успешно завершена!"
echo "Попробуйте запустить: $BINARY --help"