package service

import (
	"context"
	"errors"

	"github.com/yourusername/user-management/internal/models"
	"github.com/yourusername/user-management/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, page, pageSize int32) ([]*models.User, *Pagination, error)
}

type userService struct {
	repo repository.UserRepository
}

type Pagination struct {
	Page       int32 `json:"page"`
	PageSize   int32 `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int32 `json:"total_pages"`
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Name: req.Name,
		DOB:  req.DOB,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, ErrUserNotFound
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, ErrUserNotFound
	}

	// Check if user exists
	existingUser, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	// Update user fields
	user := &models.User{
		ID:   id,
		Name: req.Name,
		DOB:  req.DOB,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrUserNotFound
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int32) ([]*models.User, *Pagination, error) {
	// Set default values
	if page < 1 {
		page = 1
	}

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := s.repo.ListUsers(ctx, pageSize, offset)
	if err != nil {
		return nil, nil, err
	}

	// In a real application, you would get the total count from the database
	// For now, we'll just use the length of the returned users
	totalCount := int64(len(users))
	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	pagination := &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: int32(totalPages),
	}

	return users, pagination, nil
}
