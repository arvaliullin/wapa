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

lint:
	golangci-lint run ./...

down:
	docker-compose down

db:
	docker-compose up database --build --force-recreate

tests:
	go test -v ./...

env:
	- python3 -m venv env

prune:
	docker-compose down --rmi all --volumes --remove-orphans
	docker system prune -a -f
	docker volume prune -f
	docker image prune -a -f
	docker container prune -f

.PHONY: build up db tests down prune
