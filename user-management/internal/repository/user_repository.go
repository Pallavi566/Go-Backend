package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourusername/user-management/db/sqlc"
	"github.com/yourusername/user-management/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int32) ([]*models.User, error)
}

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		queries: sqlc.New(db),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	params := sqlc.CreateUserParams{
		Name: user.Name,
		Dob:  user.DOB,
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		DOB:       dbUser.Dob,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	dbUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &models.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		DOB:       dbUser.Dob,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	params := sqlc.UpdateUserParams{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.DOB,
	}

	dbUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		DOB:       dbUser.Dob,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) ListUsers(ctx context.Context, limit, offset int32) ([]*models.User, error) {
	params := sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	dbUsers, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, &models.User{
			ID:        dbUser.ID,
			Name:      dbUser.Name,
			DOB:       dbUser.Dob,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		})
	}

	return users, nil
}
