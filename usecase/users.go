package usecase

import (
	"fmt"
	"movies/model"
	"movies/repository"
)

type UsersUseCase struct {
	UsersRepo *repository.UsersRepo
}

func NewUsersUseCase(UsersRepo *repository.UsersRepo) *UsersUseCase {
	return &UsersUseCase{UsersRepo: UsersRepo}
}

// Create handles user creation logic.
func (pu *UsersUseCase) Create(user *model.Users) (*model.Users, error) {
	// Check if the email already exists in the database
	emailExists, err := pu.UsersRepo.IsEmailExists(user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return nil, fmt.Errorf("email already exists")
	}

	// Proceed to create the user if email doesn't exist
	return pu.UsersRepo.Create(user)
}

// GetUser retrieves the user from the database.
func (pu *UsersUseCase) GetUser(user *model.Users) (*model.Users, error) {
	return pu.UsersRepo.GetUser(user)
}
