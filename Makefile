.PHONY: build server-run

build: bin/server bin/client

bin/server: ./cmd/server
	go build -o bin/server ./cmd/server

bin/client: ./cmd/client
	go build -o bin/client ./cmd/client

server-run: bin/server
	bin/server

test:
	go test ./...