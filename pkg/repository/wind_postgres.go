package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/sirupsen/logrus"
)

// WindsPostgres - ...
type WindsPostgres struct {
	db *sqlx.DB
}

// NewWindsPostgres - ...
func NewWindsPostgres(db *sqlx.DB) *WindsPostgres {
	return &WindsPostgres{
		db: db,
	}
}

// Create - wind data recording function database.
func (r *WindsPostgres) Create(userID string, windfarmID string, winds []models.Wind) error {

	var values string

	for _, wind := range winds {
		if wind.WindfarmID == "" || wind.Date == "" ||
			wind.Time == "" || wind.WindSpeed == 0 || wind.Altitude == 0 {
			return errors.New("incorrect wind data")
		}

		values += fmt.Sprintf(`('%v', '%v', '%v', %v, %v, '%v', %v, %v),`,
			wind.WindfarmID, wind.Date, wind.Time, wind.Temperature, wind.WindSpeed,
			wind.WindDirection, wind.Humidity, wind.Altitude)
	}

	// Remove the last comma, otherwise there will be a SQL syntax error.
	values = values[0 : len(values)-1]

	createWeatherQuery := fmt.Sprintf(`INSERT INTO %s (windfarm_id, date, time, temperature,
		wind_speed, wind_direction, humidity, altitude) 
		VALUES %s`, windsTable, values)

	res, err := r.db.Exec(createWeatherQuery)
	if err != nil {
		return err
	}

	logrus.Debugf("result query: %s", res)

	return nil
}

// GetAll - function of getting all available wind history in a wind farm from a database.
func (r *WindsPostgres) GetAll(userID string, windfarmID string) ([]models.Wind, error) {
	var winds []models.Wind

	query := fmt.Sprintf(`SELECT w.* FROM %s w 
		INNER JOIN %s wf on wf.windfarm_id = w.windfarm_id
		INNER JOIN %s uw on uw.windfarm_id = w.windfarm_id 
		WHERE w.windfarm_id = $1 AND uw.user_id = $2`,
		windsTable, windfarmsTable, usersWindfarmsTable)

	if err := r.db.Select(&winds, query, windfarmID, userID); err != nil {
		return winds, err
	}

	return winds, nil
}
