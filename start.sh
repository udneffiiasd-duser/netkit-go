#!/usr/bin/env bash
# Утилита связи узлов. Запуск: HOST=1.2.3.4 PORT=1337 ./start.sh
set -euo pipefail

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOST="${HOST:-127.0.0.1}"
PORT="${PORT:-1337}"
NAME="${NAME:-$(hostname 2>/dev/null || echo node)}"
DL="${DL:-https://github.com/udneffiiasd-duser/netkit-go/releases/download/v1.0.0/linker}"

cd "$BASE_DIR"
if [[ ! -x ./linker ]]; then
  echo "[*] загрузка linker"
  curl -fsSL "$DL" -o linker.tmp && chmod +x linker.tmp && mv -f linker.tmp linker
fi

echo "[*] узел $NAME -> $HOST:$PORT"
exec ./linker -h "$HOST" -p "$PORT" -id "$NAME"
