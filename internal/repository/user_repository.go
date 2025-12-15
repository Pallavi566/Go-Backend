package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Pallavi566/Go-Backend/db/sqlc"
	"github.com/Pallavi566/Go-Backend/internal/models"
)

type UserRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (int64, error) {
	result, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	// First try with the generated query
	user, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		// If that fails, try with a direct query for better error messages
		var dbUser struct {
			ID        int32     `db:"id"`
			Name      string    `db:"name"`
			Dob       time.Time `db:"dob"`
			CreatedAt time.Time `db:"created_at"`
			UpdatedAt time.Time `db:"updated_at"`
		}
		err = r.db.QueryRowContext(ctx, 
			"SELECT id, name, dob, created_at, updated_at FROM users WHERE id = ?", id).
			Scan(&dbUser.ID, &dbUser.Name, &dbUser.Dob, &dbUser.CreatedAt, &dbUser.UpdatedAt)
		
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user with ID %d not found", id)
			}
			return nil, fmt.Errorf("error fetching user: %v", err)
		}
		
		return &models.User{
			ID:   int(dbUser.ID),
			Name: dbUser.Name,
			DOB:  dbUser.Dob,
		}, nil
	}
	
	return &models.User{
		ID:   int(user.ID),
		Name: user.Name,
		DOB:  user.Dob,
	}, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.User, len(users))
	for i, u := range users {
		result[i] = &models.User{
			ID:   int(u.ID),
			Name: u.Name,
			DOB:  u.Dob,
		}
	}
	return result, nil
}

func (r *UserRepository) Update(ctx context.Context, id int, name string, dob time.Time) error {
	// First, verify the user exists
	_, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		return err
	}

	// Now perform the update
	_, err = r.db.ExecContext(ctx, "UPDATE users SET name = ?, dob = ?, updated_at = NOW() WHERE id = ?",
		name, dob, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteUser(ctx, int32(id))
}

func (r *UserRepository) GetPaginated(ctx context.Context, limit, offset int) ([]*models.User, error) {
	users, err := r.queries.GetUsersPaginated(ctx, sqlc.GetUsersPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*models.User, len(users))
	for i, u := range users {
		result[i] = &models.User{
			ID:   int(u.ID),
			Name: u.Name,
			DOB:  u.Dob,
		}
	}
	return result, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

