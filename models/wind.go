package models

import (
	"errors"
)

type Wind struct {
	WeatherID     string  `json:"weather_id,omitempty" db:"wind_id"`
	WindfarmID    string  `json:"windfarm_id" db:"windfarm_id"`
	Date          string  `json:"date" db:"date"`
	Time          string  `json:"time" db:"time"`
	WindSpeed     float64 `json:"wind_speed" db:"wind_speed"`
	WindDirection float64 `json:"wind_direction" db:"wind_direction"`
	Altitude      float64 `json:"altitude" db:"altitude"`
}

//
type WindMap struct {
	Wind map[string]map[string]map[string]Wind `json:"wind"`
}

type UpdateWindInput struct {
	Date          *string  `json:"date"`
	Time          *string  `json:"time"`
	Temperature   *float64 `json:"temperature"`
	WindSpeed     *float64 `json:"wind_speed"`
	WindDirection *string  `json:"wind_direction"`
	Humidity      *float64 `json:"humidity"`
	Altitude      *float64 `json:"altitude"`
}

type ChartData struct {
	ExperementalDistribution [31]int     `json:"experemental_distribution"`
	DistributionDensity      [31]int     `json:"distribution_density"`
	WindEnergy               [31]float64 `json:"wind_energy"`
}

type TableData struct {
	Avg               map[string]map[int]float64 `json:"avg"`
	StandardDeviation map[string]map[int]float64 `json:"standard_deviation"`
	Dispersion        map[string]map[int]float64 `json:"dispersion"`
}

func (i UpdateWindInput) Validate() error {
	if i.Date == nil && i.Time == nil &&
		i.Temperature == nil && i.WindSpeed == nil &&
		i.WindDirection == nil && i.Humidity == nil &&
		i.Altitude == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
