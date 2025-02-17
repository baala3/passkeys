#! /bin/bash
set -m

# spin up postgres and redis using docker
docker compose up -d

# Generate frontend files
cd client && yarn dev &

# Start the server and watch for changes
cd server && air && fg
