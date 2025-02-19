install:
	docker compose up -d # spining up postgres and redis
	sleep 3 # waiting for services to be ready
	cd client && yarn install # installing frontend dependencies
	cd server && go run ./db/migration db init # initializing db migrations
	cd server && go run ./db/migration db migrate # running db migrations
	docker compose down # stopping docker services

run:
	set -m # enabling job control
	docker compose up -d # spining up postgres and redis using docker
	cd client && yarn dev & # generating frontend files
	cd server && air && fg # starting the server and watching for changes

dump-db:
	rm -rf ./server/db/data/** # removing all data from the db
