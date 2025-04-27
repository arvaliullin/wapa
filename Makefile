include scripts/makefiles/testing.mk

sub:
	go run github.com/arvaliullin/wapa/examples/sub

pub:
	go run github.com/arvaliullin/wapa/examples/pub

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

.PHONY: sub pub build up db tests down
