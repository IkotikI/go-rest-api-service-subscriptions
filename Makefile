MAKEFLAGS += --jobs=4

# Docker aliases
build:
	docker compose build 

up:
	docker compose up

down:
	docker compose down

# Local build and run.
# Simple 'go run' can cause OS firewall issues, because of temp files with
# random filenames.
run: 
	export CONFIG_PATH=./config/config.yaml
	make build-all
	./cmd/bin/main.exe

# Build Swagger and the App.
build-all: swag build-main

swag:
	swag init -g ./cmd/main.go

build-main:
	go build -o ./cmd/bin/main.exe ./cmd/main.go

air:
	air -c ./cmd/.air.toml

test:
	go test -v -parallel $(nproc) ./... -args -log-level=warn

clean:
	docker compose down -v
	rm -rf tmp
	rm -rf .database
# ---- DB Migrations ----
# Execute up migrations database
db-migrate-up:
	go run ./tools/migrations/goose.go up
# Execute down migrations
db-migrate-down:
	go run ./tools/migrations/goose.go down


# ---- Dev DB ----
# Run dev database in docker
dev-db-up:
	docker compose -p=dev-pg-1 -f ./tools/dev-db/docker-compose.yml up -d
# Stop dev database in docker
dev-db-down:
	docker compose -p=dev-pg-1 -f ./tools/dev-db/docker-compose.yml down
