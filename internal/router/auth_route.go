package router

import (
	"inventory-system/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterAuthRoutes sets up the routing endpoints for authentication.
func AuthRoutes(r chi.Router, authHandler handler.AuthHandler, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.With(authMiddleware).Post("/logout", authHandler.Logout)
	})
}
