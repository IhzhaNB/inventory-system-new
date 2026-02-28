package model

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleStaff      UserRole = "staff"
)

// User represents the "users" table in the database.
type User struct {
	BaseModel
	Name         string   `json:"name" db:"name"`
	Email        string   `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"`
	Role         UserRole `json:"role" db:"role"`
}
