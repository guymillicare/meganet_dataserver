# Makefile
build:
	@echo "Building the project..."
	go build -o bin/sportsbook cmd/server/main.go

run:
	@echo "Running the project..."
	go run cmd/server/main.go
