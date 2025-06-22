#!/bin/bash

# Проверка Docker
if ! command -v docker &> /dev/null; then
    echo "Ошибка: Docker не установлен"
    exit 1
fi

# Установка
sudo curl -L https://github.com/yourname/docker-cleaner/releases/latest/download/docker-cleaner-linux-amd64 -o /usr/local/bin/docker-cleaner
sudo chmod +x /usr/local/bin/docker-cleaner

echo "Установка завершена. Используйте: docker-cleaner --help"