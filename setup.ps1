# PowerShell setup script for Windows

Write-Host "🚀 Setting up User API Project..." -ForegroundColor Green

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go is not installed. Please install Go 1.21 or higher." -ForegroundColor Red
    exit 1
}

# Check if SQLC is installed
if (-not (Get-Command sqlc -ErrorAction SilentlyContinue)) {
    Write-Host "⚠️  SQLC is not installed. Installing..." -ForegroundColor Yellow
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
}

# Download dependencies
Write-Host "📦 Downloading Go dependencies..." -ForegroundColor Cyan
go mod download
go mod tidy

# Generate SQLC code
Write-Host "🔧 Generating SQLC code..." -ForegroundColor Cyan
sqlc generate

Write-Host "✅ Setup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:"
Write-Host "1. Set up your database (see README.md)"
Write-Host "2. Create a .env file with your database credentials"
Write-Host "3. Run 'go run cmd/server/main.go' to start the server"


