build-arm:
	GOARCH=arm64 GOOS=linux go build -o out/arm/bin/runner github.com/arvaliullin/wapa/cmd/runner
.PHONY: build-arm
