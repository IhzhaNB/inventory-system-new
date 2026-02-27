package router

import (
	"inventory-system/internal/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterAuthRoutes sets up the routing endpoints for authentication.
func AuthRoutes(r chi.Router, authHandler handler.AuthHandler) {
	r.Post("/login", authHandler.Login)
	// Nanti kalau ada register atau forgot-password, tinggal tambahin di sini:
	// r.Post("/register", authHandler.Register)
}
