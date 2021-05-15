package service

import (
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

type ModelsService struct {
	repo repository.Models
}

func NewModelsService(repo repository.Models) *ModelsService {
	return &ModelsService{repo: repo}
}

// Create - ...
func (s *ModelsService) Create(userID string, windfarmID string, model models.Model) (string, error) {
	return s.repo.Create(userID, windfarmID, model)
}

// GetAll - ...
func (s *ModelsService) GetAll(userID, windfarmID string) ([]models.Model, error) {
	return s.repo.GetAll(userID, windfarmID)
}

// GetMapData - ...
func (s *ModelsService) GetMapData(userID, windfarmID string, modelID string) (models.ModelMap, error) {
	model, err := s.repo.GetByID(userID, windfarmID, modelID)

	if err != nil {
		return models.ModelMap{}, err
	}

	modelMap := models.ModelMap{}

	for _, turbine := range model.Turbines {
		for _, production := range turbine.Productions {

			if modelMap.Production == nil {
				modelMap.Production = map[string]map[string]map[string]models.Production{}
			}
			if modelMap.Production[production.Date] == nil {
				modelMap.Production[production.Date] = map[string]map[string]models.Production{}
			}
			if modelMap.Production[production.Date][production.Time] == nil {
				modelMap.Production[production.Date][production.Time] = map[string]models.Production{}
			}
			modelMap.Production[production.Date][production.Time][turbine.ID] = *production
		}
	}

	return modelMap, err
}

// GetByID - ...
func (s *ModelsService) GetByID(userID, windfarmID, modelID string) (models.Model, error) {
	return s.repo.GetByID(userID, windfarmID, modelID)
}

// Delete - ...
func (s *ModelsService) Delete(userID, windfarmID, modelID string) error {
	return s.repo.Delete(userID, windfarmID)
}

// Update - ...
// func (s *ModelsService) Update(userID, windfarmID, modelID string, input models.UpdateModelInput) error {
// 	return s.repo.Update(userID, windfarmID, modelID, input)
// }
