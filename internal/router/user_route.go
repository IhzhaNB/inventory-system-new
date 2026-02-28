package router

import (
	"net/http"

	"inventory-system/internal/handler"

	"github.com/go-chi/chi/v5"
)

// UserRoutes sets up the routing endpoints for user management operations.
// It applies the authentication middleware to protect all user-related endpoints.
func UserRoutes(r chi.Router, userHandler handler.UserHandler, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		// Apply the authentication middleware to all endpoints within this group.
		r.Use(authMiddleware)

		// Create a new user (POST /api/v1/users)
		r.Post("/", userHandler.CreateUser)
	})
}
