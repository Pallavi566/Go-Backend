package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yourusername/user-management/config"
	"github.com/yourusername/user-management/db"
	"github.com/yourusername/user-management/internal/handler"
	"github.com/yourusername/user-management/internal/logger"
	"github.com/yourusername/user-management/internal/repository"
	"github.com/yourusername/user-management/internal/routes"
	"github.com/yourusername/user-management/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)\n
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize logger
	logg, err := logger.NewLogger("user-management", "development")
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logg.Sync()

	// Initialize database
	dbConfig := &db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	database, err := db.NewDatabase(dbConfig)
	if err != nil {
		logg.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(database.Pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// Setup routes
	routes.SetupRoutes(app, userHandler, logg)

	// Start server in a goroutine
	go func() {
		logg.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logg.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logg.Info("Shutting down server...")

	// Create a deadline to wait for
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close database connection
	if err := database.Close(); err != nil {
		logg.Error("Error closing database connection", zap.Error(err))
	}

	// Shutdown Fiber app
	if err := app.Shutdown(); err != nil {
		logg.Error("Error shutting down server", zap.Error(err))
	}

	logg.Info("Server gracefully stopped")
}

// errorHandler handles errors returned from HTTP handlers
func errorHandler(c *fiber.Ctx, err error) error {
	// Default status code
	code := fiber.StatusInternalServerError

	// Check if it's a fiber.Error
	e, ok := err.(*fiber.Error)
	if ok {
		code = e.Code
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
