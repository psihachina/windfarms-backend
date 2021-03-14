package repository

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/sirupsen/logrus"
)

// TurbinePostgres - ..
type TurbinePostgres struct {
	db *sqlx.DB
}

// NewTurbinePostgres - ..
func NewTurbinePostgres(db *sqlx.DB) *TurbinePostgres {
	return &TurbinePostgres{db: db}
}

// Create - function of recording a turbine and its power in the database
func (r *TurbinePostgres) Create(userID string, turbine models.Turbine, outputs models.Outputs) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var id string

	createTurbineQuery := fmt.Sprintf(`
		INSERT INTO %s (user_id, turbine_name, maximum_power, max_wind_speed,
		min_wind_speed, rotor_diameter, tower_height, number_blades, annual_turbine_maintenance) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING turbine_id`, turbinesTable)
	row := tx.QueryRow(createTurbineQuery, userID, turbine.TurbineName,
		turbine.MaximumPower, turbine.MaxWindSpeed, turbine.MinWindSpeed,
		turbine.RotorDiameter, turbine.TowerHeight, turbine.Blades,
		turbine.AnnualTurbineMaintenance)
	fmt.Println(id)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return "", err
	}

	var values string
	var speed int = 2

	for _, output := range outputs.Outputs {
		values += fmt.Sprintf(`('%v', '%v', '%v'),`,
			id, speed, output.Production)

		speed = speed + 2
	}

	// Remove the last comma, otherwise there will be a SQL syntax error.
	values = values[0 : len(values)-1]

	createOutputsQuery := fmt.Sprintf("INSERT INTO %s (turbine_id, speed, production) VALUES %s", outputsTable, values)
	_, err = tx.Exec(createOutputsQuery)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return id, tx.Commit()
}

// GetAll - function of getting all user turbines from the database.
func (r *TurbinePostgres) GetAll(userID string) ([]models.Turbine, error) {
	var turbines []models.Turbine

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, turbinesTable)
	err := r.db.Select(&turbines, query, userID)

	for index, _ := range turbines {
		var outputs []models.Output
		query = fmt.Sprintf(`SELECT * FROM %s WHERE turbine_id = $1`, outputsTable)
		err = r.db.Select(&outputs, query, turbines[index].TurbineID)
		turbines[index].Outputs = outputs
	}

	return turbines, err
}

// GetByID - function of obtaining user turbine by ID from the database.
func (r *TurbinePostgres) GetByID(userID string, turbineID string) (models.Turbine, error) {
	var turbine models.Turbine
	var outputs []models.Output

	query := fmt.Sprintf(`SELECT * FROM %s  WHERE user_id = $1 AND turbine_id = $2`, turbinesTable)
	err := r.db.Get(&turbine, query, userID, turbineID)

	query = fmt.Sprintf(`SELECT o.* FROM %s o JOIN %s t ON t.turbine_id = o.turbine_id 
						  WHERE t.user_id = $1 AND o.turbine_id = $2`, outputsTable, turbinesTable)

	err = r.db.Select(&outputs, query, userID, turbineID)

	turbine.Outputs = outputs

	return turbine, err
}

// Delete - function to delete user turbine by ID from the database.
func (r *TurbinePostgres) Delete(userID string, turbineID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND turbine_id = $2`, turbinesTable)
	_, err := r.db.Exec(query, userID, turbineID)

	return err
}

// Update - function of updating user turbine data by ID in the database.
func (r *TurbinePostgres) Update(userID string, turbineID string, input models.UpdateTurbineInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	v := reflect.ValueOf(input)

	for i := 0; i < v.NumField(); i++ {
		if reflect.Indirect(v.Field(i)).IsValid() {
			setValues = append(setValues, fmt.Sprintf("%s=$%d", v.Type().Field(i).Tag.Get("json"), argID))

			args = append(args, (v.Field(i).Elem().Interface()))
			argID++
		}
	}
	fmt.Println("setValues:", setValues)

	setQuery := strings.Join(setValues, ", ")
	fmt.Println("setQuery:", setQuery)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE turbine_id=$%d AND user_id=$%d",
		turbinesTable, setQuery, argID, argID+1)
	args = append(args, turbineID, userID)
	fmt.Println(args)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}
