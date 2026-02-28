package middleware

import (
	"context"
	"net/http"
	"strings"

	"inventory-system/internal/repository"
	"inventory-system/pkg/utils"

	"github.com/google/uuid"
)

// Create a custom type for Context Keys to prevent collisions with other packages.
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UserRoleKey ContextKey = "user_role"
)

// Authenticate verifies the validity of the UUID token against the database.
func Authenticate(sessionRepo repository.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. Extract the token from the authorization header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.Error(w, r, http.StatusUnauthorized, "Missing authorization header", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.Error(w, r, http.StatusUnauthorized, "Invalid authorization format", nil)
				return
			}

			// 2. Parse the token string into a UUID.
			sessionID, err := uuid.Parse(parts[1])
			if err != nil {
				utils.Error(w, r, http.StatusUnauthorized, "Invalid token format", nil)
				return
			}

			// 3. Query the database to check if the session is valid.
			session, err := sessionRepo.GetValid(r.Context(), sessionID)
			if err != nil {
				// If the session is not found, expired, or revoked, reject the request.
				utils.Error(w, r, http.StatusUnauthorized, "Token is expired or invalid", nil)
				return
			}

			// 4. If valid, store the UserID and Role into the Request Context.
			// This allows subsequent endpoints (e.g., /items) to identify the authenticated user.
			ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, session.Role)

			// Proceed to the next handler with the populated context.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
