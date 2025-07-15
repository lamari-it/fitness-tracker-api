.PHONY: build run test clean deps

# Build the application
build:
	go build -o bin/fitflow-api main.go

# Run the application
run:
	go run main.go

# Run with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Create database
createdb:
	createdb fitflow

# Drop database
dropdb:
	dropdb fitflow

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/fitflow-api main.go