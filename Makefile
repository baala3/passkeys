install:
	docker compose up -d # spining up postgres and redis
	sleep 1 # waiting for services to be ready
	cd client && yarn install # installing frontend dependencies
	cd server && go run ./db/migration db init # initializing db migrations
	cd server && go run ./db/migration db migrate # running db migrations
	docker compose down # stopping docker services

run:
	set -m # enabling job control
	docker compose up -d # spining up postgres and redis using docker
	cd client && yarn watch & # generating frontend files & watching for changes
	cd server && air && fg # starting the server and watching for changes

dump-db:
	docker compose down -v # removing all data from the db
	rm -rf ./server/data/services/postgres/ # removing all data from the db
	rm -rf ./server/data/services/redis/ # removing all data from the db

dump-data:
	rm -rf ./server/data/services/postgres/ # removing all data from the db
	rm -rf ./server/data/services/redis/ # removing all data from the db

wire:
	cd server && wire
