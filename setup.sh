#! /bin/bash

# Install Frontend de
cd client && yarn install

# setup Database
docker compose up -d
cd ../server && go run ./db/migration db init && go run ./db/migration db migrate 
docker compose down
