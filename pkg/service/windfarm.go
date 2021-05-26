package service

import (
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

// WindfarmsService - ...
type WindfarmsService struct {
	repo repository.Windfarms
}

// NewWindfarmsService - ...
func NewWindfarmsService(repo repository.Windfarms) *WindfarmsService {

	return &WindfarmsService{repo: repo}
}

// Create - ...
func (s *WindfarmsService) Create(userID string, windfarm models.Windfarm) (string, error) {

	return s.repo.Create(userID, windfarm)
}

// GetAll - ...
func (s *WindfarmsService) GetAll(userID string) ([]models.Windfarm, error) {
	return s.repo.GetAll(userID)
}

// GetByID - ...
func (s *WindfarmsService) GetByID(userID, windfarmID string) (models.Windfarm, error) {
	return s.repo.GetByID(userID, windfarmID)
}

// Delete - ...
func (s *WindfarmsService) Delete(userID, windfarmID string) error {
	return s.repo.Delete(userID, windfarmID)
}

// Update - ...
func (s *WindfarmsService) Update(userID string, itemID string, input models.UpdateWindfarmInput) error {
	return s.repo.Update(userID, itemID, input)
}
