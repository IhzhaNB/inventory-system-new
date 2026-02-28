package router

import (
	_ "inventory-system/docs"
	"inventory-system/internal/handler"
	customMiddleware "inventory-system/internal/middleware"
	"inventory-system/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Setup initializes the main chi router, attaches middlewares, and registers all sub-routes.
func SetupRoute(handlers *handler.Handler, repos *repository.Repository) *chi.Mux {
	r := chi.NewRouter()

	// Standard Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authMiddleware := customMiddleware.Authenticate(repos.Session)
	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API Versioning Group
	r.Route("/api/v1", func(r chi.Router) {

		// Register module routes here
		AuthRoutes(r, handlers.Auth, authMiddleware)
		UserRoutes(r, handlers.User, authMiddleware)

	})

	return r
}
