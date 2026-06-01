#!/bin/sh
set -eu

echo "[ENTRYPOINT] Starting"

LITESTREAM_CONFIG="${LITESTREAM_CONFIG:-/etc/litestream.yml}"
MAIN_BIN="${MAIN_BIN:-/usr/local/bin/main}"

DATABASE_URL="${DATABASE_URL:-}"
if [ -z "${DATABASE_URL}" ]; then
  echo "[ENTRYPOINT] DATABASE_URL is empty; starting app without litestream"
  exec "${MAIN_BIN}"
fi

# Extract the file path from DATABASE_URL (supports "file:/path", "file:///path", or plain "/path")
DB_PATH="${DATABASE_URL}"
case "${DB_PATH}" in
  file://*) DB_PATH="${DB_PATH#file://}" ;;
  file:*) DB_PATH="${DB_PATH#file:}" ;;
esac

# If litestream config isn't present, just run the app.
if [ ! -f "${LITESTREAM_CONFIG}" ]; then
  echo "[ENTRYPOINT] No litestream config at ${LITESTREAM_CONFIG}; starting app without litestream"
  exec "${MAIN_BIN}"
fi

# If replica env isn't configured, don't block app startup on litestream.
if [ -z "${DB_REPLICA_URL:-}" ] || [ -z "${R2_BUCKET:-}" ] || [ -z "${R2_ACCESS_KEY_ID:-}" ] || [ -z "${R2_SECRET_ACCESS_KEY:-}" ]; then
  echo "[ENTRYPOINT] Litestream replica env not fully set; starting app without litestream"
  echo "[ENTRYPOINT] Replica env status: DB_REPLICA_URL=$([ -n "${DB_REPLICA_URL:-}" ] && echo set || echo missing), R2_BUCKET=$([ -n "${R2_BUCKET:-}" ] && echo set || echo missing), R2_ACCESS_KEY_ID=$([ -n "${R2_ACCESS_KEY_ID:-}" ] && echo set || echo missing), R2_SECRET_ACCESS_KEY=$([ -n "${R2_SECRET_ACCESS_KEY:-}" ] && echo set || echo missing)"
  exec "${MAIN_BIN}"
fi

echo "[ENTRYPOINT] Starting Litestream (config=${LITESTREAM_CONFIG}, db=${DB_PATH})"
echo "[ENTRYPOINT] Litestream replica target: endpoint=${DB_REPLICA_URL}, bucket=${R2_BUCKET}, access_key_id_len=${#R2_ACCESS_KEY_ID}, secret_access_key_len=${#R2_SECRET_ACCESS_KEY}"
litestream version || true

# Restore the database if it doesn't exist
if [ ! -f "${DB_PATH}" ]; then
  echo "[ENTRYPOINT] Litestream restore (if replica exists)"
  litestream restore -config "${LITESTREAM_CONFIG}" -if-replica-exists -v "${DB_PATH}"
fi

echo "[ENTRYPOINT] Litestream replicate (exec=${MAIN_BIN})"
export LITESTREAM_ACTIVE=true
exec litestream replicate -config "${LITESTREAM_CONFIG}" -exec "${MAIN_BIN}"
