package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/Pallavi566/Go-Backend/config"
	"github.com/Pallavi566/Go-Backend/internal/handler"
	"github.com/Pallavi566/Go-Backend/internal/logger"
	"github.com/Pallavi566/Go-Backend/internal/repository"
	"github.com/Pallavi566/Go-Backend/internal/routes"
	"github.com/Pallavi566/Go-Backend/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create a context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Log.Info("Starting User API Server...")

	// Initialize database connection
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Test database connection
	if err := db.PingContext(ctx); err != nil {
		logger.Log.Fatal("Failed to ping database", zap.Error(err))
	}

	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(*userRepo)
	userHandler := handler.NewUserHandler(userService, logger.Log)

	// Initialize Fiber app with context support
	app := fiber.New(fiber.Config{
		AppName:               "User API",
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(requestid.New())

	// Add context to each request
	app.Use(func(c *fiber.Ctx) error {
		// Set the context with timeout for each request
		reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Set the context in the locals
		c.Locals("ctx", reqCtx)

		// Override c.UserContext to return our context
		c.SetUserContext(reqCtx)

		return c.Next()
	})

	// Setup routes
	routes.SetupRoutes(app, userHandler, logger.Log)

	// Start server in a goroutine
	serverShutdown := make(chan error, 1)
	go func() {
		addr := "0.0.0.0:" + cfg.GetPort()
		logger.Log.Info("Server starting", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			serverShutdown <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case <-ctx.Done():
		logger.Log.Info("Received shutdown signal")
	case err := <-serverShutdown:
		logger.Log.Error("Server error", zap.Error(err))
	}

	// Create a deadline for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	if err := db.Close(); err != nil {
		logger.Log.Error("Error closing database connection", zap.Error(err))
	}

	logger.Log.Info("Server gracefully stopped")
}

func errorHandler(c *fiber.Ctx, err error) error {
	logger.Log.Error("Unhandled error", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal server error",
	})
}
