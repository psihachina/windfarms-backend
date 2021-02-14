package repository

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	if err := initConfig(); err != nil {
		t.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("error loading env virables: %s", err.Error())
	}

	db, err := NewPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewAuthPostgres(db)

	tests := []struct {
		name    string
		input   models.User
		wantErr bool
	}{
		{
			name: "Ok",
			input: models.User{
				Email:    "test@test.ru",
				Password: "test",
			},
		},
		{
			name: "Empty fields(this Ok)",
			input: models.User{
				Email:    "",
				Password: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.CreateUser(tt.input)
			assert.NoError(t, err)
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	if err := initConfig(); err != nil {
		t.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("error loading env virables: %s", err.Error())
	}

	db, err := NewPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewAuthPostgres(db)

	type args struct {
		email    string
		password string
	}

	tests := []struct {
		name    string
		input   args
		wantErr bool
		want    models.User
	}{
		{
			name:  "Ok",
			input: args{"test@test.ru", "test"},
			want: models.User{
				UserID:   "610be5ec-4ab3-400e-8562-7fe736a35eb6",
				Email:    "",
				Password: "",
			},
		},
		{
			name:    "Not Found",
			input:   args{"not@test.ru", "found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetUser(tt.input.email, tt.input.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func initConfig() error {
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
