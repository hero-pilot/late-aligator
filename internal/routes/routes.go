package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/hero-pilot/late-aligator/internal/app"
)



func SetUpRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	return r
}