#!/usr/bin/env bash
# Утилита связи узлов. Запуск: ./start.sh
set -euo pipefail

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOST="${HOST:-127.0.0.1}"
PORT="${PORT:-1337}"
NAME="${NAME:-$(hostname 2>/dev/null || echo node)}"

if [[ ! -x "$BASE_DIR/linker" ]]; then
  echo "[!] файл linker не найден рядом со start.sh" >&2
  exit 1
fi

echo "[*] узел $NAME -> $HOST:$PORT"
exec "$BASE_DIR/linker" -h "$HOST" -p "$PORT" -id "$NAME"
