.PHONY: build run test clean migrate-up migrate-down sqlc

# Variables
BINARY_NAME=user-management
MIGRATION_DIR=db/migrations
DB_URL=postgresql://postgres:postgres@localhost:5432/userdb?sslmode=disable

# Build the application
build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build files
clean:
	go clean
	rm -f bin/$(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy
	go mod download

# Generate SQLC code
sqlc:
	sqlc generate

# Run database migrations
migrate-up:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up

# Rollback the latest migration
migrate-down:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" down

# Show migration status
migrate-status:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" status

# Start a PostgreSQL container
db-up:
	docker run --name user-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=userdb -p 5432:5432 -d postgres:13-alpine

# Stop and remove the PostgreSQL container
db-down:
	docker stop user-db
	docker rm user-db

# Run the application with Docker
docker-run:
	docker build -t user-management .
	docker run --name user-management --network host -p 3000:3000 --env-file .env user-management
