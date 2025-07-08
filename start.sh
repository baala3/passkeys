#!/bin/sh
set -e

# Wait for database to be ready
echo "==> Waiting for database to be ready..."
while ! nc -z postgres 5432; do
  sleep 1
done

# Wait for Redis to be ready
echo "==> Waiting for Redis to be ready..."
while ! nc -z redis 6379; do
  sleep 1
done

echo "==> Database and Redis are ready!"

# .env file is not loaded in the container, 
# so we need to export the variables manually to pass DB init & migration commands
export DB_HOST=postgres
export DB_PORT=5432
export DB_USER=myuser
export DB_PASSWORD=mypassword
export DB_NAME=mydb

echo "==> Starting DB init..."
go run ./db/migration db init

echo "==> Running DB migrations..."
go run ./db/migration db migrate

echo "==> Starting backend app..."
./main
