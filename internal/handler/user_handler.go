package handler

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/Pallavi566/Go-Backend/internal/models"
	"github.com/Pallavi566/Go-Backend/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	service  *service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", req.DOB); err != nil {
		h.logger.Error("Invalid date format", zap.String("dob", req.DOB), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Expected YYYY-MM-DD",
		})
	}

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	h.logger.Info("User created", zap.Int("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		h.logger.Error("User not found", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()
	users, err := h.service.GetAllUsers(ctx)
	if err != nil {
		h.logger.Error("Failed to get users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	return c.JSON(users)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate date format if DOB is provided
	if req.DOB != "" {
		if _, err := time.Parse("2006-01-02", req.DOB); err != nil {
			h.logger.Error("Invalid date format", zap.String("dob", req.DOB), zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected YYYY-MM-DD",
			})
		}
	}

	user, err := h.service.UpdateUser(ctx, id, req)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	h.logger.Info("User updated", zap.Int("user_id", id))
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.service.DeleteUser(ctx, id); err != nil {
		h.logger.Error("Failed to delete user", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	h.logger.Info("User deleted", zap.Int("user_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) GetUsersPaginated(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var params models.PaginationParams
	if err := c.QueryParser(&params); err != nil {
		h.logger.Error("Failed to parse query parameters", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if err := h.validate.Struct(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := h.service.GetUsersPaginated(ctx, params.Page, params.Limit)
	if err != nil {
		h.logger.Error("Failed to get paginated users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(result)
}

