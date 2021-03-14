package models

//User ...
type User struct {
	UserID   string `json:"-" db:"user_id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
