
helloworld:
	go run github.com/arvaliullin/wapa-composer/examples/helloworld
sub:
	go run github.com/arvaliullin/wapa-composer/examples/sub
pub:
	go run github.com/arvaliullin/wapa-composer/examples/pub

swag:
	swag init -g cmd/composer/main.go

build: swag
	go build -o out/bin/composer github.com/arvaliullin/wapa-composer/cmd/composer

up:
	docker-compose up --build --force-recreate

db:
	docker-compose up database --build --force-recreate

tests:
	go test -v ./...

.PHONY: helloworld sub pub build up db tests
