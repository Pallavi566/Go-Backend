.PHONY: run build test migrate sqlc docker-up docker-down clean

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Generate SQLC code
sqlc:
	sqlc generate

# Run migrations (manual - you'll need to run this against your DB)
migrate:
	@echo "Run migrations manually using: mysql -u user -p userdb < db/migrations/001_create_users_table.sql"

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Clean generated files
clean:
	rm -rf bin/
	rm -rf db/sqlc/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Full setup
setup: deps sqlc
	@echo "Setup complete! Run 'make run' to start the server"


