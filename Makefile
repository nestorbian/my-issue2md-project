.PHONY: test build web

test:
	go test ./...

build:
	go build -o bin/issue2md ./cmd/cli

web:
	go build -o bin/issue2md-web ./cmd/web
