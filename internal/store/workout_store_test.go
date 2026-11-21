package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", "postgres://shitman:shitman@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("Openning test db: %v", err)
	}
	t.Cleanup(func() {db.Close()})
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}
	_, err = db.Exec(`TRUNCATE workout, workout_entries CASCADE`)

	if err != nil {
		t.Fatalf("truncating tables error: %v", err)
	}
	return db
}


func TestCreateWorkout(t *testing.T) {
	db := setUpTestDB(t)

	store := NewPostgresWorkoutStore(db)

	tests := []struct{
		name string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title: "Push day",
				Description: "Upper body day",
				DurationMinutes: 60,
				CaloriesBurned: 200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench press",
						Sets: 3,
						Reps: IntPtr(10),
						Weight: FloatPtr(135.5),
						Notes: "Warmup properly",
						OrderIndex: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Workout with invalid entries",
			workout: &Workout{
				Title: "Full body",
				Description: "Complete workout",
				DurationMinutes: 90,
				CaloriesBurned: 200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets: 3,
						Reps: IntPtr(60),
						Notes: "Keep form",
						OrderIndex: 1,
					},
					{
						ExerciseName: "Squats",
						Sets: 4,
						Reps: IntPtr(12),
						DurationSeconds: IntPtr(50) ,
						Weight: FloatPtr(156.8),
						Notes: "Shake dat ass",
						OrderIndex: 2,
					},
				},
			},
			wantErr: true,
			},
		}

	for _, testCase:= range tests{
		t.Run(testCase.name , func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(testCase.workout)
			if testCase.wantErr{
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.workout.Title, createdWorkout.Title)
			assert.Equal(t, testCase.workout.Description, createdWorkout.Description)
			assert.Equal(t, testCase.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(testCase.workout.Entries), len(createdWorkout.Entries))

			for i, _ := range retrieved.Entries{
				assert.Equal(t, testCase.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, testCase.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, testCase.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}
func FloatPtr(i float64) *float64 {
	return &i
}