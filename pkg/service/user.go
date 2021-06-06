package service

import (
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

// UsersService - ...
type UsersService struct {
	repo repository.Users
}

// NewUsersService - ...
func NewUsersService(repo repository.Users) *UsersService {

	return &UsersService{repo: repo}
}

// GetAll - ...
func (s *UsersService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

// Delete - ...
func (s *UsersService) Delete(email string) error {
	return s.repo.Delete(email)
}

// Update - ...
func (s *UsersService) Confirm(email string) error {
	return s.repo.Confirm(email)
}
