package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Pallavi566/Go-Backend/internal/handler"
	"github.com/Pallavi566/Go-Backend/internal/middleware"
	"go.uber.org/zap"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, logger *zap.Logger) {
	// Apply global middleware
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware(logger))

	// Health check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1 routes
	api := app.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.Post("/", userHandler.CreateUser)
			users.Get("/", userHandler.GetUsersPaginated) // Paginated by default
			users.Get("/all", userHandler.GetAllUsers)    // Get all without pagination
			users.Get("/:id", userHandler.GetUserByID)
			users.Put("/:id", userHandler.UpdateUser)
			users.Delete("/:id", userHandler.DeleteUser)
		}
	}
}


