package model

// User represents the "users" table in the database.
type User struct {
	BaseModel
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"` // Hidden from JSON responses
	Role         string `json:"role" db:"role"`
}
