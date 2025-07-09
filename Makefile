install:
	docker compose -f compose-dev.yaml up -d # spining up postgres and redis
	sleep 2 # waiting for services to be ready
	cd client && yarn install # installing frontend dependencies
	cd server && go run ./db/migration db init # initializing db migrations
	cd server && go run ./db/migration db migrate # running db migrations
	docker compose -f compose-dev.yaml down # stopping docker services

run:
	set -m # enabling job control
	docker compose -f compose-dev.yaml up -d # spining up postgres and redis using docker
	cd client && yarn watch & # generating frontend files & watching for changes
	cd server && air && fg # starting the server and watching for changes

dump-db:
	docker compose -f compose-dev.yaml down -v # removing all data from the db
	rm -rf ./server/data/services/postgres/ # removing all data from the db
	rm -rf ./server/data/services/redis/ # removing all data from the db

wire:
	cd server && wire
