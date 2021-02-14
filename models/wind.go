package models

import "errors"

type Wind struct {
	WeatherID     string  `json:"weather_id,omitempty" db:"wind_id"`
	WindfarmID    string  `json:"windfarm_id" db:"windfarm_id"`
	Date          string  `json:"date" db:"date"`
	Time          string  `json:"time" db:"time"`
	Temperature   float64 `json:"temperature" db:"temperature"`
	WindSpeed     float64 `json:"wind_speed" db:"wind_speed"`
	WindDirection string  `json:"wind_direction" db:"wind_direction"`
	Humidity      float64 `json:"humidity" db:"humidity"`
	Altitude      float64 `json:"altitude" db:"altitude"`
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

func (i UpdateWindInput) Validate() error {
	if i.Date == nil && i.Time == nil &&
		i.Temperature == nil && i.WindSpeed == nil &&
		i.WindDirection == nil && i.Humidity == nil &&
		i.Altitude == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
