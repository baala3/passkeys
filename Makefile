install:
	docker compose up -d # spin up postgres and redis
	sleep 3 # wait for services to be ready
	cd client && yarn install && cd .. # Install Frontend dependencies
	cd server && go run ./db/migration db init && go run ./db/migration db migrate && cd .. # Run migrations
	docker compose down # Stop docker services

run:
	set -m # Enable job control
	docker compose up -d # spin up postgres and redis using docker
	cd client && yarn dev & # Generate frontend files
	cd server && air && fg # Start the server and watch for changes

