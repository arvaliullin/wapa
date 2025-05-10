include scripts/makefiles/testing.mk
include scripts/makefiles/arm.mk

swag:
	swag init -g cmd/composer/main.go

build: swag
	go build -o out/bin/composer github.com/arvaliullin/wapa/cmd/composer

up:
	docker-compose up --build --force-recreate

down:
	docker-compose down -v --rmi all

db:
	docker-compose up database --build --force-recreate

tests:
	go test -v ./...


DB_USER=postgres
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres

export PGUSER=$(DB_USER)
export PGPASSWORD=$(DB_PASSWORD)
export PGHOST=$(DB_HOST)
export PGPORT=$(DB_PORT)
export PGDATABASE=$(DB_NAME)

export-db:
	docker-compose exec -T database pg_dump -U $(DB_USER) $(DB_NAME) > out/db_dump.sql

import-db:
	cat out/db_dump.sql | docker-compose exec -T database psql -U $(DB_USER) $(DB_NAME)

.PHONY: sub pub build up db tests down export-db import-db
