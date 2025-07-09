#!/bin/sh
set -e

echo "==> Starting DB init..."
go run ./db/migration db init

echo "==> Running DB migrations..."
go run ./db/migration db migrate

echo "==> Starting backend app..."
./main
