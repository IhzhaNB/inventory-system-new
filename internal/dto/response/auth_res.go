package response

import "github.com/google/uuid"

// UserResponse represents the safe, public user data returned to the client.
// Notice it does not include sensitive fields like PasswordHash.
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

// AuthResponse defines the JSON response sent back after a successful login.
// It combines the JWT token and the sanitized UserResponse.
type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}
