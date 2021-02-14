package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable          = "users"
	windfarmsTable      = "windfarms"
	usersWindfarmsTable = "users_windfarms"
	windsTable          = "winds"
)

// Config ...
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

//NewPostgresDB ...
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
