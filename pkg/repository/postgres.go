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
	turbinesTable       = "turbines"
	productions         = "productions"
	outputsTable        = "outputs"
	modelTable          = "models"
	trubinesModelsTabel = "turbines_models"
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

func (c *Config) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", c.Host, c.Port, c.Username, c.DBName, c.Password, c.SSLMode)
}

//NewPostgresDB ...
func NewPostgresDB(cfg string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
