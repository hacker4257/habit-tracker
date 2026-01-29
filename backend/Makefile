.PHONY: build run test clean dev

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run: build
	./bin/server

# Run in development mode
dev:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./cmd/server

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/server-mac ./cmd/server
