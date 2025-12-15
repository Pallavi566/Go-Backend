package service

import (
	"context"
	"time"

	"github.com/Pallavi566/Go-Backend/internal/models"
	"github.com/Pallavi566/Go-Backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return nil, err
	}

	age := calculateAge(dob)
	return &models.UserResponse{
		ID:   int(id),
		Name: req.Name,
		DOB:  dob.Format("2006-01-02"),
		Age:  &age,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	age := calculateAge(user.DOB)
	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format("2006-01-02"),
		Age:  &age,
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []models.UserResponse
	for _, u := range users {
		age := calculateAge(u.DOB)
		response = append(response, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.DOB.Format("2006-01-02"),
			Age:  &age,
		})
	}

	return response, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req models.UpdateUserRequest) (*models.UserResponse, error) {
    // Get the existing user to preserve fields not being updated
    existingUser, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Use existing values if not provided in request
    name := existingUser.Name
    if req.Name != "" {
        name = req.Name
    }

    // Parse DOB if provided, otherwise use existing DOB
    var dob time.Time
    if req.DOB != "" {
        dob, err = time.Parse("2006-01-02", req.DOB)
        if err != nil {
            return nil, err
        }
    } else {
        dob = existingUser.DOB
    }

    // Update the user
    err = s.repo.Update(ctx, id, name, dob)
    if err != nil {
        return nil, err
    }

    // Get the updated user to ensure we return the latest data
    updatedUser, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    age := calculateAge(updatedUser.DOB)
    return &models.UserResponse{
        ID:   updatedUser.ID,
        Name: updatedUser.Name,
        DOB:  updatedUser.DOB.Format("2006-01-02"),
        Age:  &age,
    }, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) GetUsersPaginated(ctx context.Context, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	users, err := s.repo.GetPaginated(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	var response []models.UserResponse
	for _, u := range users {
		age := calculateAge(u.DOB)
		response = append(response, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.DOB.Format("2006-01-02"),
			Age:  &age,
		})
	}

	return &models.PaginatedResponse{
		Data:       response,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (int(total) + limit - 1) / limit,
	}, nil
}

// calculateAge calculates age from date of birth
func calculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// If birthday hasn't occurred yet this year, subtract one year
	if now.YearDay() < dob.YearDay() {
		years--
	}

	return years
}
