package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	DOB       time.Time `json:"dob" validate:"required,validateDOB"`
	Age       int       `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name string    `json:"name" validate:"required,min=2,max=100"`
	DOB  time.Time `json:"dob" validate:"required,validateDOB"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name string    `json:"name" validate:"required,min=2,max=100"`
	DOB  time.Time `json:"dob" validate:"required,validateDOB"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.DOB.Format("2006-01-02"),
		Age:  u.CalculateAge(),
	}
}

// CalculateAge calculates the user's age based on their date of birth
func (u *User) CalculateAge() int {
	now := time.Now()
	year, month, day := now.Date()
	dobYear, dobMonth, dobDay := u.DOB.Date()

	age := year - dobYear

	// If birthday hasn't occurred yet this year, subtract one year
	if month < dobMonth || (month == dobMonth && day < dobDay) {
		age--
	}

	return age
}

// Validate validates the user struct
func (u *User) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("validateDOB", validateDOB)
	return validate.Struct(u)
}

// validateDOB is a custom validator for date of birth
func validateDOB(fl validator.FieldLevel) bool {
	dob, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Check if the date is in the past
	if dob.After(time.Now()) {
		return false
	}

	// Check if the age is reasonable (e.g., less than 150 years)
	age := time.Since(dob).Hours() / (24 * 365)
	return age <= 150
}
