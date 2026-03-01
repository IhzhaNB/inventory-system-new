package request

// LoginRequest defines the JSON payload expected from the frontend during login.
// We use validation tags later here to ensure data integrity.
type LoginRequest struct {
	Email    string `json:"email" example:"admin@gmail.com"`
	Password string `json:"password" example:"password123"`
}
