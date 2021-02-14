package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/psihachina/windfarms-backend/models"
)

// AuthPostgres - ...
type AuthPostgres struct {
	db *sqlx.DB
}

// NewAuthPostgres - ...
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateUser - ...
func (r *AuthPostgres) CreateUser(user models.User) error {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (email, password_hash) values ($1, $2) RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

// GetUser - ...
func (r *AuthPostgres) GetUser(email, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	return user, err
}
