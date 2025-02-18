#! /bin/bash

# spin up postgres and redis
docker compose up -d

# wait for services to be ready
sleep 3
# Install Frontend dependencies
cd client && yarn install && cd ..

# Run migrations
cd server && go run ./db/migration db init && go run ./db/migration db migrate && cd ..

# Stop docker services
docker compose down
