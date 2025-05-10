include scripts/makefiles/testing.mk
include scripts/makefiles/arm.mk
include scripts/makefiles/plot.mk
include scripts/makefiles/db.mk

swag:
	swag init -g cmd/composer/main.go

build: swag
	go build -o out/bin/composer github.com/arvaliullin/wapa/cmd/composer

up:
	docker-compose up --build --force-recreate

down:
	docker-compose down

db:
	docker-compose up database --build --force-recreate

tests:
	go test -v ./...

env:
	- python3 -m venv env

.PHONY: build up db tests down
