#!/bin/bash

echo "🚀 Setting up User API Project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check if SQLC is installed
if ! command -v sqlc &> /dev/null; then
    echo "⚠️  SQLC is not installed. Installing..."
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy

# Generate SQLC code
echo "🔧 Generating SQLC code..."
sqlc generate

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Set up your database (see README.md)"
echo "2. Create a .env file with your database credentials"
echo "3. Run 'make run' to start the server"


