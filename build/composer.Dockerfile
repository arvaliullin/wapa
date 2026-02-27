FROM golang:1.26.0-trixie

WORKDIR /repo

COPY go.mod ./

RUN go mod download

COPY . .

RUN mkdir -p /opt/wapa/composer
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go build -o /opt/wapa/composer/bin/composer github.com/arvaliullin/wapa/cmd/composer

COPY configs/composer/config.yaml /etc/wapa/composer/config.yaml

CMD ["/opt/wapa/composer/bin/composer", "-config=/etc/wapa/composer/config.yaml"]
