# Network Monitor Makefile

.PHONY: all build clean run-client run-server build-client build-server build-web

all: build

build: build-client build-server build-web

# Build both client and server
build-client:
	go build -o bin/client ./cmd/client

build-server:
	go build -o bin/server ./cmd/server

build-web:
	cd web && npm install && npm run build

# Run the client
run-client:
	go run ./cmd/client/main.go

# Run the server
run-server:
	go run ./cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin
	rm -rf web/dist

# Install dependencies
deps:
	go mod download
	cd web && npm install