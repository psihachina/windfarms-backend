package models

import "encoding/json"

type Model struct {
	ModelID    string         `json:"model_id" db:"model_id"`
	ModelName  string         `json:"model_name" db:"model_name"`
	WindfarmID string         `json:"windfarm_id" db:"windfarm_id"`
	Turbines   []TurbineModel `json:"turbines,omitempty" db:"turbines"`
}

type TurbineModel struct {
	TurbineModelID string      `json:"turbines_models_id" db:"turbines_models_id"`
	ID             string      `json:"id," db:"id"`
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
	ICUF           float64 `json:"icuf,omitempty" db:"icuf"`
	Value          float64 `json:"value,omitempty" db:"value"`
	TurbineModelID string  `json:"turbines_models_id" db:"turbines_models_id"`
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
