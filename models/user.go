package models

//User ...
type User struct {
	UserID       string `json:"-" db:"user_id"`
	Email        string `json:"email" db:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	RegisteredAt string `json:"registered_at" db:"registered_at"`
	AdminConfirm bool   `json:"admin_confirm" db:"admin_confirm"`
	EmailConfirm bool   `json:"email_confirm" db:"email_confirm"`
}
