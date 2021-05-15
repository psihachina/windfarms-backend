package service

import (
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

// Authorization - ...
type Authorization interface {
	CreateUser(user models.User) error
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (string, error)
}

// Models - ...
type Models interface {
	Create(userID string, windfarmID string, model models.Model) (string, error)
	GetAll(userID, windfarmID string) ([]models.Model, error)
	GetMapData(userID, windfarmID string, modelID string) (models.ModelMap, error)
	GetByID(userID, windfarmID, modelID string) (models.Model, error)
	Delete(userID, windfarmID, modelID string) error
	//Update(userID, windfarmID, modelID string, inputModel models.UpdateModelInput) error
}

// Windfarms - ...
type Windfarms interface {
	Create(userID string, windfarm models.Windfarm) (string, error)
	GetAll(userID string) ([]models.Windfarm, error)
	GetByID(userID string, windfarmID string) (models.Windfarm, error)
	Delete(userID string, windfarmID string) error
	Update(userID string, windfarmID string, input models.UpdateWindfarmInput) error
}

// Turbines - ...
type Turbines interface {
	Create(userID string, turbine models.Turbine, outputs models.Outputs) (string, error)
	GetAll(userID string) ([]models.Turbine, error)
	GetByID(userID, turbineID string) (models.Turbine, error)
	Delete(userID string, turbineID string) error
	Update(userID string, turbineID string, input models.UpdateTurbineInput) error
}

// Winds - ...
type Winds interface {
	Create(userID string, windfarmID string) error
	GetAll(userID, windfarmID string) ([]models.Wind, error)
	GetWindForChart(userID, windfarmID, from, to string, height int) (models.ChartData, error)
	GetWindForTable(userID, windfarmID string) (models.TableData, error)
}

// Service - ...
type Service struct {
	Authorization
	Windfarms
	Winds
	Turbines
	Models
}

// NewService - ...
func NewService(repos *repository.Repository) *Service {
	windfarmsService := NewWindfarmsService(repos.Windfarms)
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Windfarms:     NewWindfarmsService(repos.Windfarms),
		Winds:         NewWindsService(repos.Winds, windfarmsService),
		Turbines:      NewTurbinesService(repos.Turbines),
		Models:        NewModelsService(repos.Models),
	}
}
