package models

import "encoding/json"

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
	Productions    Productions `json:"productions" db:"productions"`
}

type Production struct {
	ProductionID   string  `json:"production_id,omitempty" db:"production_id"`
	Time           string  `json:"time,omitempty" db:"time"`
	Date           string  `json:"date,omitempty" db:"date"`
	ICUF           float64 `json:"icuf" db:"icuf"`
	Value          float64 `json:"value,omitempty" db:"value"`
	WindSpeed      float64 `json:"wind_speed,omitempty" db:"wind_speed"`
	TurbineModelID string  `json:"turbines_models_id" db:"turbines_models_id"`
	Altitude       float64 `json:"altitude" db:"altitude"`
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
