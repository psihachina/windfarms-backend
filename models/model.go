package models

import (
	"encoding/json"
	"errors"
)

type Model struct {
	ModelID          string         `json:"model_id" db:"model_id"`
	WindfarmAltitude float64        `json:"windfarm_altitude"`
	ModelName        string         `json:"model_name" db:"model_name"`
	WindfarmID       string         `json:"windfarm_id" db:"windfarm_id"`
	Turbines         []TurbineModel `json:"turbines,omitempty" db:"turbines"`
}

type TurbineModel struct {
	TurbineModelID string      `json:"turbines_models_id" db:"turbines_models_id"`
	TurbineName    string      `json:"turbine_name," db:"turbine_name"`
	Latitude       float64     `json:"latitude," db:"latitude"`
	Longitude      float64     `json:"longitude," db:"longitude"`
	ModelID        string      `json:"model_id" db:"model_id"`
	X              float64     `json:"x" db:"x"`
	Y              float64     `json:"y" db:"y"`
	Z              float64     `json:"z" db:"z"`
	Productions    Productions `json:"productions" db:"productions"`
}

type Production struct {
	ProductionID     string  `json:"production_id,omitempty" db:"production_id"`
	Time             string  `json:"time,omitempty" db:"time"`
	Date             string  `json:"date,omitempty" db:"date"`
	ICUF             float64 `json:"icuf" db:"icuf"`
	Value            float64 `json:"value,omitempty" db:"value"`
	WindSpeed        float64 `json:"wind_speed,omitempty" db:"wind_speed"`
	TurbineModelID   string  `json:"turbines_models_id" db:"turbines_models_id"`
	Altitude         float64 `json:"altitude" db:"altitude"`
	WindDirection    float64 `json:"wind_direction" db:"wind_direction"`
	Shading          float64 `json:"shading" db:"shading"`
	SpeedWithShading float64 `json:"speed_with_shading" db:"speed_with_shading"`
}

type ModelMap struct {
	Production map[string]map[string]map[string]Production `json:"production"`
}

type Productions []*Production

func (ls *Productions) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, ls)
	}
	return nil
}

type UpdateModelInput struct {
	ModelName *string `json:"model_name"`
}

func (i UpdateModelInput) Validate() error {
	if i.ModelName == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
