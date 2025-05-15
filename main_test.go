package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func clearTableData(db *sql.DB) error {
	clearDataQuery := `
		DELETE FROM sets;
		DELETE FROM workouts;
	`
	_, err := db.Exec(clearDataQuery)
	return err
}


func TestCreateWorkoutHandler(t *testing.T) {
	initDB("workouts_test.db") 
	clearTableData(db) 
	defer db.Close()

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

	assert.Equal(t, http.StatusCreated, w.Code)

	rows, err := db.Query("SELECT id, date, exercise FROM workouts")
	assert.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 1, count, "Expected one workout in the database")
}
