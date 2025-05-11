package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/mattn/go-sqlite3"
)

// Added cleanup logic to clear the database before each test
func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		exercise TEXT
	);
		CREATE TABLE IF NOT EXISTS sets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workout_id INTEGER,
		weight INTEGER,
		reps INTEGER,
		FOREIGN KEY(workout_id) REFERENCES workouts(id)
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		panic(err)
	}

	// Clear any existing data
	db.Exec("DELETE FROM workouts")
	db.Exec("DELETE FROM sets")

	return db
}

// Updated test to invoke initDB from main.go
func TestCreateWorkoutHandler(t *testing.T) {
	// Setup test database
	initDB() // Use the shared initDB function from main.go
	defer db.Close()

	// Create a test HTTP server
	reqBody := map[string]interface{}{
		"date":     "2025-05-11",
		"exercise": "Squat",
		"sets": []map[string]interface{}{
			{"weight": 100, "reps": 5},
			{"weight": 120, "reps": 3},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/workouts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(createWorkoutHandler)
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify data in the database
	rows, err := db.Query("SELECT id, date, exercise FROM workouts")
	assert.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 1, count, "Expected one workout in the database")
}
