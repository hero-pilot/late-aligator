package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setUpTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "postgres://shitman:shitman@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("Openning test db: %w", err)
	}
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %w", err)
	}
	_, err = db.Exec(`TRUNCATE workout, workout_entries CASCADE`)

	if err != nil {
		t.Fatalf("truncating tables error: %w", err)
	}
	return db
}
