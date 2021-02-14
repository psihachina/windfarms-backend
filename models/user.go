package models

//User ...
type User struct {
	UserID   string `json:"-" db:"user_id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
