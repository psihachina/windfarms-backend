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

// Windfarms - ...
type Windfarms interface {
	Create(userID string, item models.Windfarm) (string, error)
	GetAll(userID string) ([]models.Windfarm, error)
	GetByID(userID, windfarmID string) (models.Windfarm, error)
	Delete(userID string, windfarmID string) error
	Update(userID string, windfarmID string, input models.UpdateWindfarmInput) error
}

// Models - ...
type Models interface {
	Create(userID string, windfarmID string, model models.Model) (string, error)
	GetAll(userID, windfarmID string) ([]models.Model, error)
	GetByID(userID, windfarmID, modelID string) (models.Model, error)
	Delete(userID string, windfarmID string) error
	//Update(userID, windfarmID, modelID string, input models.UpdateModelInput) error
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
	Create(userID string, windfarmID string, winds []models.Wind) error
	GetAll(userID string, windfarmID string) ([]models.Wind, error)
	GetWindForChart(userID, windfarmID, from, to string, height int) ([]models.Wind, error)
}

// Repository - ...
type Repository struct {
	Authorization
	Windfarms
	Winds
	Turbines
	Models
}

// NewRepository - ...
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Windfarms:     NewWindfarmPostgres(db),
		Winds:         NewWindsPostgres(db),
		Turbines:      NewTurbinePostgres(db),
		Models:        NewModelPostgres(db),
	}
}
