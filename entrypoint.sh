#!/bin/sh
set -x

echo "[ENTRYPOINT] Starting Litestream"

# Extract the file path from DATABASE_URL
DB_PATH=$(echo $DATABASE_URL | sed 's/^file://')

# Restore the database if it doesn't exist
if [ ! -f $DB_PATH ]; then
    echo "Starting Litestream [Restore]"
    litestream restore -if-replica-exists -v $DB_PATH
fi

echo "Starting Litestream [Replicate]"

# Start Litestream in the background
litestream replicate -exec "main"

echo "[ENTRYPOINT] [DONE] Started Litestream"