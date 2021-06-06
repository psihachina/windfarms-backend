package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
)

// UserPostgres - ..
type UserPostgres struct {
	db *sqlx.DB
}

// NewUserPostgres - ..
func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// GetAll - ...
func (r *UserPostgres) GetAll() ([]models.User, error) {
	var model []models.User

	query := fmt.Sprintf(`SELECT email, registered_at, admin_confirm, email_confirm FROM %s`, usersTable)
	err := r.db.Select(&model, query)

	return model, err
}

// Delete - ...
func (r *UserPostgres) Delete(email string) error {
	query := fmt.Sprintf(`DELETE FROM %s 
						WHERE email = $1`, usersTable)
	_, err := r.db.Exec(query, email)

	return err
}

//Update - ...
func (r *UserPostgres) Confirm(email string) error {
	query := fmt.Sprintf(`UPDATE %s SET admin_confirm = $1 
						WHERE email = $2`, usersTable)
	_, err := r.db.Exec(query, true, email)

	return err
}
