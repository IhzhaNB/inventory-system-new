package response

// AuthResponse defines the JSON response sent back after a successful login.
// It combines the JWT token and the sanitized UserResponse.
type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}
