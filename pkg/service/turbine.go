package service

import (
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

// TurbinesService - ...
type TurbinesService struct {
	repo repository.Turbines
}

// NewTurbinesService - ...
func NewTurbinesService(repo repository.Turbines) *TurbinesService {

	return &TurbinesService{repo: repo}
}

// Create - ...
func (s *TurbinesService) Create(userID string, turbine models.Turbine, outputs models.Outputs) (string, error) {
	return s.repo.Create(userID, turbine, outputs)
}

// GetAll - ...
func (s *TurbinesService) GetAll(userID string) ([]models.Turbine, error) {
	return s.repo.GetAll(userID)
}

// GetByID - ...
func (s *TurbinesService) GetByID(userID, turbineID string) (models.Turbine, error) {
	return s.repo.GetByID(userID, turbineID)
}

// Delete - ...
func (s *TurbinesService) Delete(userID, turbineID string) error {
	return s.repo.Delete(userID, turbineID)
}

// Update - ...
func (s *TurbinesService) Update(userID string, turbineID string, input models.UpdateTurbineInput) error {
	return s.repo.Update(userID, turbineID, input)
}
