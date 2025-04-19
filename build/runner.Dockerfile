FROM golang:latest

WORKDIR /repo

COPY go.mod ./

RUN go mod download

COPY . .

RUN mkdir -p /opt/wapa/runner
RUN go build -o /opt/wapa/runner/bin/runner github.com/arvaliullin/wapa/cmd/runner

COPY configs/runner/config.yaml /etc/wapa/runner/config.yaml

CMD ["/opt/wapa/runner/bin/runner", "-config=/etc/wapa/runner/config.yaml"]
