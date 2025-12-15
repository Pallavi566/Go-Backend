package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger is a middleware that logs HTTP requests
func Logger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID from context
		reqID := c.GetRespHeader("X-Request-ID")

		// Log request details
		fields := []zap.Field{
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("ip", c.IP()),
			zap.Duration("latency", latency),
		}

		// Add error field if there was an error
		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		// Log the request
		logger.Info("HTTP request", fields...)

		return err
	}
}

// RequestID is a middleware that adds a unique request ID to the response headers
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateRequestID()
		}

		// Add request ID to response headers
		c.Set("X-Request-ID", reqID)

		// Add request ID to context
		ctx := context.WithValue(c.Context(), "request_id", reqID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// In a real application, you might want to use a UUID or similar
	return time.Now().Format("20060102150405") + "-abc123"
}
