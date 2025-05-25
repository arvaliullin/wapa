FROM golang:latest
RUN apt update && \
    apt install -y ca-certificates curl unzip hyperfine && \
    rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:${PATH}"


WORKDIR /repo

COPY go.mod ./

RUN go mod download

COPY . .

RUN mkdir -p /opt/wapa/runner
RUN mkdir -p /opt/wapa/bin
RUN mkdir -p /opt/wapa/scripts

RUN go build -o /opt/wapa/runner/bin/runner github.com/arvaliullin/wapa/cmd/runner

COPY configs/runner/config.yaml /etc/wapa/runner/config.yaml
COPY scripts/runner /opt/wapa/scripts

CMD ["/opt/wapa/runner/bin/runner", "-config=/etc/wapa/runner/config.yaml"]
