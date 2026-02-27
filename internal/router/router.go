package router

import (
	"inventory-system/internal/handler"

	_ "inventory-system/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Setup initializes the main chi router, attaches middlewares, and registers all sub-routes.
func SetupRoute(handler *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Standard Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Poin ke file json docs
	))

	// API Versioning Group
	r.Route("/api/v1", func(r chi.Router) {

		// Register module routes here
		AuthRoutes(r, handler.Auth)

		// Nanti lu tinggal nambahin modul lain dengan gampang:
		// RegisterItemRoutes(r, handlers.Item)
		// RegisterSaleRoutes(r, handlers.Sale)

		// Kalau butuh rute yang diprotect JWT, bisa di-group di dalam sini juga nanti
	})

	return r
}
