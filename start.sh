#!/bin/sh
set -e

echo "Running database migrations..."
# Source app.env if it exists (for baked-in secrets from build)
if [ -f /app/app.env ]; then
    set -a
    . /app/app.env
    set +a
fi

/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Starting the server..."
exec "$@"
