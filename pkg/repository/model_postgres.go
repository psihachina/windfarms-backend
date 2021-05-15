package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
)

// ModelPostgres - ..
type ModelPostgres struct {
	db *sqlx.DB
}

// NewModelPostgres - ..
func NewModelPostgres(db *sqlx.DB) *ModelPostgres {
	return &ModelPostgres{db: db}
}

// Create - ..
func (r *ModelPostgres) Create(userID string, windfarmID string, model models.Model) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var id string
	if err != nil {
		return "", err
	}

	createModelQuery := fmt.Sprintf(`
		INSERT INTO %s (model_name, windfarm_id) 
		VALUES ($1, $2) 
		RETURNING model_id`, modelTable)
	row := tx.QueryRow(createModelQuery, model.ModelName, windfarmID)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	var values string

	for _, turbine := range model.Turbines {
		fmt.Println(turbine.ID)
		values += fmt.Sprintf(`('%v', '%v', '%v', '%v', '%v'),`,
			id, turbine.TurbineName, turbine.Latitude, turbine.Latitude, turbine.ID)
	}

	// Remove the last comma, otherwise there will be a SQL syntax error.
	values = values[0 : len(values)-1]

	createOutputsQuery := fmt.Sprintf("INSERT INTO %s (model_id, turbine_name, latitude, longitude, id) VALUES %s", trubinesModelsTabel, values)
	_, err = tx.Exec(createOutputsQuery)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return id, tx.Commit()
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

	for i, m := range model {
		query = fmt.Sprintf(`SELECT tm.*, JSON_AGG(TO_JSON(p.*)) as productions FROM %s tm LEFT JOIN %s p ON p.turbines_models_id = tm.turbines_models_id
							WHERE model_id = $1 GROUP BY tm.turbines_models_id`, trubinesModelsTabel, productions)
		err = r.db.Select(&model[i].Turbines, query, m.ModelID)

		if err != nil {
			return model, err
		}
	}

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

	query = fmt.Sprintf(`SELECT tm.*, JSON_AGG(TO_JSON(p.*)) as productions FROM %s tm LEFT JOIN %s p ON p.turbines_models_id = tm.turbines_models_id
							WHERE model_id = $1 GROUP BY tm.turbines_models_id`, trubinesModelsTabel, productions)
	err = r.db.Select(&model.Turbines, query, modelID)

	if err != nil {
		return model, err
	}

	return model, err
}

// Delete - ...
func (r *ModelPostgres) Delete(userID string, windfarmID string) error {
	query := fmt.Sprintf(`DELETE FROM %s wf 
						USING %s uw WHERE wf.windfarm_id = uw.windfarm_id AND uw.user_id = $1 AND uw.windfarm_id = $2`, modelTable, trubinesModelsTabel)
	_, err := r.db.Exec(query, userID, windfarmID)

	return err
}

// Update - ...
// func (r *ModelPostgres) Update(userID, windfarmID, modelID string, input models.UpdateModelInput) error {
// 	setValues := make([]string, 0)
// 	args := make([]interface{}, 0)
// 	argID := 1

// 	v := reflect.ValueOf(input)

// 	for i := 0; i < v.NumField(); i++ {
// 		if reflect.Indirect(v.Field(i)).IsValid() {
// 			setValues = append(setValues, fmt.Sprintf("%s=$%d", v.Type().Field(i).Tag.Get("json"), argID))
// 			args = append(args, (v.Field(i).Elem().Interface().(string)))
// 			vl := v.Field(i).Elem().Interface().(string)
// 			fmt.Println(vl)
// 			argID++
// 		}
// 	}

// 	setQuery := strings.Join(setValues, ", ")

// 	query := fmt.Sprintf("UPDATE %s wf SET %s FROM %s uw WHERE wf.windfarm_id = uw.windfarm_id AND uw.windfarm_id=$%d AND uw.user_id=$%d",
// 		modelTable, setQuery, modelTrubinesTabel, argID, argID+1)
// 	args = append(args, windfarmID, userID)
// 	fmt.Println(args)

// 	logrus.Debugf("updateQuery: %s", query)
// 	logrus.Debugf("args: %s", args)
// 	_, err := r.db.Exec(query, args...)
// 	return err
// }
