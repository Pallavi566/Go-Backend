package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/user-management/internal/handler"
	"github.com/yourusername/user-management/internal/middleware"
	"go.uber.org/zap"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, logger *zap.Logger) {
	// Middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(logger))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	{
		users.Post("/", userHandler.CreateUser)
		users.Get("/", userHandler.ListUsers)
		users.Get("/:id", userHandler.GetUser)
		users.Put("/:id", userHandler.UpdateUser)
		users.Delete("/:id", userHandler.DeleteUser)
	}

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not Found",
		})
	})
}
