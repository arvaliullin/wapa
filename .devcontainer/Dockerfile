FROM mcr.microsoft.com/devcontainers/go:1.24

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN apt-get update && apt-get install -y docker-ce-cli docker-compose

RUN apt-get install -y git
