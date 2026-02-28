package response

import (
	"inventory-system/internal/model"

	"github.com/google/uuid"
)

// UserResponse represents the safe, public user data returned to the client.
// Notice it does not include sensitive fields like PasswordHash.
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func ToUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
}
