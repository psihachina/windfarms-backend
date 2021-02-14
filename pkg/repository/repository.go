package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
)

// Authorization - ...
type Authorization interface {
	CreateUser(user models.User) error
	GetUser(username, password string) (models.User, error)
}

type Windfarms interface {
	Create(userID string, item models.Windfarm) (string, error)
	GetAll(userID string) ([]models.Windfarm, error)
	GetByID(userID, windfarmID string) (models.Windfarm, error)
	Delete(userID string, windfarmID string) error
	Update(userID string, windfarmID string, input models.UpdateWindfarmInput) error
}

// Winds - ...
type Winds interface {
	Create(userID string, windfarmID string, winds []models.Wind) error
	GetAll(userID string, windfarmID string) ([]models.Wind, error)
}

// Repository - ...
type Repository struct {
	Authorization
	Windfarms
	Winds
}

// NewRepository - ...
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Windfarms:     NewWindfarmPostgres(db),
		Winds:         NewWindsPostgres(db),
	}
}
