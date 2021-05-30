package repository

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/sirupsen/logrus"
)

// ModelPostgres - ..
type ModelPostgres struct {
	db *sqlx.DB
}

// NewModelPostgres - ..
func NewModelPostgres(db *sqlx.DB) *ModelPostgres {
	return &ModelPostgres{db: db}
}

//CreateModel - ...
func (r *ModelPostgres) CreateModel(userID string, windfarmID string, model models.Model) (string, error) {
	var id string

	createModelQuery := fmt.Sprintf(`
		INSERT INTO %s (model_name, windfarm_id) 
		VALUES ($1, $2) 
		RETURNING model_id`, modelTable)

	row := r.db.QueryRow(createModelQuery, model.ModelName, windfarmID)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

// GenerateModel - ..
func (r *ModelPostgres) GenerateModel(userID, windfarmID, modelID string, model models.Model) (string, error) {
	var wg sync.WaitGroup

	for _, turbine := range model.Turbines {
		t := turbine
		wg.Add(1)

		go func(t models.TurbineModel) {
			defer wg.Done()

			createTurbineModelQuery := fmt.Sprintf("INSERT INTO %s (turbines_models_id ,model_id, turbine_name, latitude, longitude, x, y, z) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (turbines_models_id) DO NOTHING", trubinesModelsTabel)
			_, err := r.db.Exec(createTurbineModelQuery, t.TurbineModelID, modelID, t.TurbineName, t.Latitude, t.Longitude, t.X, t.Y, t.Z)
			if err != nil {
				log.Fatal(err)
			}

			var values string

			for _, production := range t.Productions {
				values += fmt.Sprintf(`('%v', '%v', '%v', '%v', '%v', '%v', '%v'),`,
					production.Value, production.ICUF, production.WindSpeed, production.Date, production.Time, t.TurbineModelID, production.Altitude)
			}

			values = values[0 : len(values)-1]
			createProductionsQuery := fmt.Sprintf(`INSERT INTO %s (value , icuf, wind_speed, date, time, turbines_models_id, altitude) 
													VALUES %s `, productions, values)

			_, err = r.db.Exec(createProductionsQuery)
			if err != nil {
				log.Fatal(err)
			}
		}(t)
	}

	wg.Wait()

	return modelID, nil
}

// GetAll - ...
func (r *ModelPostgres) GetAll(userID, windfarmID string) ([]models.Model, error) {
	var model []models.Model

	query := fmt.Sprintf(`SELECT m.* FROM %s m 
						INNER JOIN %s wf on wf.windfarm_id = m.windfarm_id 
						INNER JOIN %s uw on uw.windfarm_id = wf.windfarm_id
						WHERE uw.user_id = $1 AND uw.windfarm_id = $2
						`, modelTable, windfarmsTable, usersWindfarmsTable)
	err := r.db.Select(&model, query, userID, windfarmID)

	return model, err
}

// GetByID - ...
func (r *ModelPostgres) GetByID(userID, windfarmID, modelID string) (models.Model, error) {
	var model models.Model

	query := fmt.Sprintf(`SELECT m.* FROM %s m 
						INNER JOIN %s wf on wf.windfarm_id = m.windfarm_id 
						INNER JOIN %s uw on uw.windfarm_id = wf.windfarm_id
						WHERE uw.user_id = $1 AND uw.windfarm_id = $2 AND m.model_id = $3
						`, modelTable, windfarmsTable, usersWindfarmsTable)
	err := r.db.Get(&model, query, userID, windfarmID, modelID)
	if err != nil {
		return model, err
	}

	query = fmt.Sprintf(`SELECT tm.* FROM %s tm
							WHERE model_id = $1 GROUP BY tm.turbines_models_id`, trubinesModelsTabel)
	err = r.db.Select(&model.Turbines, query, modelID)

	if err != nil {
		return model, err
	}

	return model, err
}

// GetByID - ...
func (r *ModelPostgres) GetByIDMap(userID, windfarmID, modelID string) (models.Model, error) {
	var model models.Model

	query := fmt.Sprintf(`SELECT m.* FROM %s m 
						INNER JOIN %s wf on wf.windfarm_id = m.windfarm_id 
						INNER JOIN %s uw on uw.windfarm_id = wf.windfarm_id
						WHERE uw.user_id = $1 AND uw.windfarm_id = $2 AND m.model_id = $3
						`, modelTable, windfarmsTable, usersWindfarmsTable)
	err := r.db.Get(&model, query, userID, windfarmID, modelID)
	if err != nil {
		return model, err
	}

	query = fmt.Sprintf(`SELECT tm.*,JSON_AGG(TO_JSON(p.*)) as productions FROM %s tm LEFT JOIN %s p ON p.turbines_models_id = tm.turbines_models_id
							WHERE model_id = $1 GROUP BY tm.turbines_models_id`, trubinesModelsTabel, productions)
	err = r.db.Select(&model.Turbines, query, modelID)

	if err != nil {
		return model, err
	}

	return model, err
}

// Delete - ...
func (r *ModelPostgres) Delete(userID, windfarmID, modelID string) error {
	query := fmt.Sprintf(`DELETE FROM %s m 
						USING %s uw WHERE uw.windfarm_id = m.windfarm_id AND uw.user_id = $1 AND m.windfarm_id = $2 AND m.model_id = $3`, modelTable, usersWindfarmsTable)
	_, err := r.db.Exec(query, userID, windfarmID, modelID)

	return err
}

// DeleteTurbine - ...
func (r *ModelPostgres) DeleteTurbine(modelID, modelTrubineID string) error {
	query := fmt.Sprintf(`DELETE FROM %s tm 
						USING %s m WHERE tm.model_id = m.model_id AND m.model_id = $1 
						AND tm.turbines_models_id = $2 `, trubinesModelsTabel, modelTable)
	_, err := r.db.Exec(query, modelID, modelTrubineID)

	return err
}

//Update - ...
func (r *ModelPostgres) Update(userID, windfarmID, modelID string, input models.UpdateModelInput) error {
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

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s m SET %s FROM %s uw WHERE m.windfarm_id = uw.windfarm_id AND uw.windfarm_id=$%d AND uw.user_id=$%d AND m.model_id=$%d",
		modelTable, setQuery, usersWindfarmsTable, argID, argID+1, argID+2)
	args = append(args, windfarmID, userID, modelID)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}
