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

// Windfarms - ...
type Windfarms interface {
	Create(userID string, windfarm models.Windfarm) (string, error)
	GetAll(userID string) ([]models.Windfarm, error)
	GetByID(userID string, windfarmID string) (models.Windfarm, error)
	Delete(userID string, windfarmID string) error
	Update(userID string, windfarmID string, input models.UpdateWindfarmInput) error
}

// Winds - ...
type Winds interface {
	Create(userID string, windfarmID string) error
	GetAll(userID, windfarmID string) ([]models.Wind, error)
}

// Service - ...
type Service struct {
	Authorization
	Windfarms
	Winds
}

// NewService - ...
func NewService(repos *repository.Repository) *Service {
	windfarmsService := NewWindfarmsService(repos.Windfarms)
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Windfarms:     NewWindfarmsService(repos.Windfarms),
		Winds:         NewWindsService(repos.Winds, windfarmsService),
	}
}
