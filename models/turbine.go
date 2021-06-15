package models

import (
	"errors"
)

type Turbine struct {
	TurbineID     string   `json:"turbineID,omitempty" db:"turbine_id"`
	UserID        string   `json:"userID" db:"user_id"`
	TurbineName   string   `json:"turbine_name" db:"turbine_name" binding:"required"`
	MaximumPower  float64  `json:"maximum_power" db:"maximum_power" binding:"required"`
	MaxWindSpeed  float64  `json:"max_wind_speed" db:"max_wind_speed" binding:"required"`
	MinWindSpeed  float64  `json:"min_wind_speed" db:"min_wind_speed" binding:"required"`
	Blades        int      `json:"number_blades" db:"number_blades" binding:"required"`
	TowerHeight   float64  `json:"tower_height" db:"tower_height" binding:"required"`
	RotorDiameter float64  `json:"rotor_diameter" db:"rotor_diameter" binding:"required"`
	Outputs       []Output `json:"outputs" db:"outputs" binding:"required"`
}

// Outputs struct only to get data from the body
type Outputs struct {
	Outputs []Output
}

type Output struct {
	OutputID   string  `json:"outputID,omitempty" db:"output_id"`
	TurbineID  string  `json:"turbineID" db:"turbine_id"`
	Speed      int     `json:"speed" db:"speed"`
	Production float64 `json:"production" db:"production"`
}

type UpdateTurbineInput struct {
	TurbineName              *string  `json:"turbine_name"`
	MaximumPower             *float64 `json:"maximum_power"`
	MaxWindSpeed             *float64 `json:"max_wind_speed"`
	MinWindSpeed             *float64 `json:"min_wind_speed"`
	RotorDiameter            *float64 `json:"rotor_diameter"`
	AnnualTurbineMaintenance *float64 `json:"annual_turbine_maintenance"`
}

func (i UpdateTurbineInput) Validate() error {
	if i.TurbineName == nil && i.MaximumPower == nil &&
		i.RotorDiameter == nil &&
		i.MaxWindSpeed == nil && i.MinWindSpeed == nil &&
		i.AnnualTurbineMaintenance == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
