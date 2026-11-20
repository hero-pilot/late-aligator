package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hero-pilot/late-aligator/internal/store"
	"github.com/hero-pilot/late-aligator/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: Decoding createWorkout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request sent"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: CreateWorkout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})

}

func (wh *WorkoutHandler) HandleUpdatingWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		ID              *int                 `json:"id"`
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("ERROR: Decoding: %v", err)
		utils.WriteJson(w, http.StatusBadGateway, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkoutByID(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: UpdateWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	err = wh.workoutStore.DeleteWorkoutByID(workoutID)
	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: DeleteWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}
	if err != nil {
		wh.logger.Printf("ERROR: DeleteWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusNoContent, utils.Envelope{})
}
