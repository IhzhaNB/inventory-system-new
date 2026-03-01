package router

import (
	"net/http"

	"inventory-system/internal/handler"
	customMiddleware "inventory-system/internal/middleware"
	"inventory-system/internal/model"

	"github.com/go-chi/chi/v5"
)

// UserRoutes sets up the routing endpoints for user management operations.
func UserRoutes(r chi.Router, userHandler handler.UserHandler, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		// 1. PRIMARY GATE: Authentication (Check if user is logged in via valid UUID token)
		r.Use(authMiddleware)

		// 2. SECONDARY GATE: Authorization (Check if user role is super_admin or admin)
		// Variadic function allows us to pass multiple allowed roles easily.
		r.Use(customMiddleware.RequireRole(
			string(model.RoleSuperAdmin),
			string(model.RoleAdmin),
		))

		// 3. ENDPOINTS: Only accessible if BOTH gates above are passed.
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUsers)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})
}
