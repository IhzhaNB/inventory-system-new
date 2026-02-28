package router

import (
	"inventory-system/internal/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterAuthRoutes sets up the routing endpoints for authentication.
func AuthRoutes(r chi.Router, authHandler handler.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		// r.Post("/register", authHandler.Register)
	})
}
