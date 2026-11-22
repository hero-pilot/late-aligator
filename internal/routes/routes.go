package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/hero-pilot/late-aligator/internal/app"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdatingWorkoutByID)
	r.Delete("/workouts/{id}" , app.WorkoutHandler.HandleDeleteWorkoutByID)

	r.Post("/users", app.UserHandler.HandleRegisterUser)
	return r
}
