package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hero-pilot/late-aligator/internal/api"
)

type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will  go here

	//our handlers will go here
	workouteHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger: logger,
		WorkoutHandler: workouteHandler,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "status is available\n")
	a.Logger.Println("Get: /HealthCheck, 200")
}
