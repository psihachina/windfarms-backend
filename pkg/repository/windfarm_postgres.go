package repository

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/sirupsen/logrus"
)

// WindfarmPostgres - ..
type WindfarmPostgres struct {
	db *sqlx.DB
}

// NewWindfarmPostgres - ..
func NewWindfarmPostgres(db *sqlx.DB) *WindfarmPostgres {
	return &WindfarmPostgres{db: db}
}

// Create - ..
func (r *WindfarmPostgres) Create(userID string, windfarm models.Windfarm) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var id string

	createWindfarmQuery := fmt.Sprintf(`
		INSERT INTO %s (windfarm_name, polygon, longitude,
		latitude, capacity, range_to_city, range_to_road, range_to_city_line,
		city_longitude, city_latitude, description, polygon_radius
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING windfarm_id`, windfarmsTable)
	row := tx.QueryRow(createWindfarmQuery, windfarm.WindfarmName,
		windfarm.PolygonJSON, windfarm.Longitude,
		windfarm.Latitude, windfarm.Capacity, windfarm.RangeToCity,
		windfarm.RangeToRoad, windfarm.RangeToCityLine,
		windfarm.CityLongitude, windfarm.CityLatitude, windfarm.Description, windfarm.PolygonRadius)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	createUsersWindfarmsQuery := fmt.Sprintf("INSERT INTO %s (user_id, windfarm_id) VALUES($1, $2) RETURNING users_windfarms_id", usersWindfarmsTable)
	_, err = tx.Exec(createUsersWindfarmsQuery, userID, id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return id, tx.Commit()
}

// GetAll - ...
func (r *WindfarmPostgres) GetAll(userID string) ([]models.Windfarm, error) {
	var windfarms []models.Windfarm

	query := fmt.Sprintf(`SELECT wf.* FROM %s wf INNER JOIN %s uw on wf.windfarm_id = uw.windfarm_id WHERE uw.user_id = $1`, windfarmsTable, usersWindfarmsTable)
	err := r.db.Select(&windfarms, query, userID)

	return windfarms, err
}

// GetByID - ...
func (r *WindfarmPostgres) GetByID(userID string, windfarmID string) (models.Windfarm, error) {
	var windfarm models.Windfarm

	query := fmt.Sprintf(`SELECT wf.* FROM %s wf 
						INNER JOIN %s uw on wf.windfarm_id = uw.windfarm_id WHERE uw.user_id = $1 AND uw.windfarm_id = $2`, windfarmsTable, usersWindfarmsTable)
	err := r.db.Get(&windfarm, query, userID, windfarmID)

	return windfarm, err
}

// Delete - ...
func (r *WindfarmPostgres) Delete(userID string, windfarmID string) error {
	query := fmt.Sprintf(`DELETE FROM %s wf 
						USING %s uw WHERE wf.windfarm_id = uw.windfarm_id AND uw.user_id = $1 AND uw.windfarm_id = $2`, windfarmsTable, usersWindfarmsTable)
	_, err := r.db.Exec(query, userID, windfarmID)

	return err
}

// Update - ...
func (r *WindfarmPostgres) Update(userID string, windfarmID string, input models.UpdateWindfarmInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	v := reflect.ValueOf(input)

	for i := 0; i < v.NumField(); i++ {
		if reflect.Indirect(v.Field(i)).IsValid() {
			setValues = append(setValues, fmt.Sprintf("%s=$%d", v.Type().Field(i).Tag.Get("json"), argID))
			args = append(args, (v.Field(i).Elem().Interface().(string)))
			vl := v.Field(i).Elem().Interface().(string)
			fmt.Println(vl)
			argID++
		}
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s wf SET %s FROM %s uw WHERE wf.windfarm_id = uw.windfarm_id AND uw.windfarm_id=$%d AND uw.user_id=$%d",
		windfarmsTable, setQuery, usersWindfarmsTable, argID, argID+1)
	args = append(args, windfarmID, userID)
	fmt.Println(args)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}
