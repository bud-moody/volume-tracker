package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var logger *slog.Logger

func openDB(fileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTablesIfNotPresent(db *sql.DB) error {
	createTablesQuery := `
		CREATE TABLE IF NOT EXISTS workouts (
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
		);
	`
	_, err := db.Exec(createTablesQuery)
	return err
}

func initDB(dbName string) {
	var err error
	db, err = openDB(dbName)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err := createTablesIfNotPresent(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}

func initLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

type Set struct {
	Weight int `json:"weight"`
	Reps   int `json:"reps"`
}

type Workout struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Exercise string `json:"exercise"`
	Sets     []Set  `json:"sets"`
}

func createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))

	var workout Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO workouts (date, exercise) VALUES (?, ?)", workout.Date, workout.Exercise)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	workoutID, _ := result.LastInsertId()
	for _, set := range workout.Sets {
		_, err := db.Exec("INSERT INTO sets (workout_id, weight, reps) VALUES (?, ?, ?)", workoutID, set.Weight, set.Reps)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func getWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))

	rows, err := db.Query("SELECT id, date, exercise FROM workouts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workouts []Workout
	for rows.Next() {
		var workout Workout
		if err := rows.Scan(&workout.ID, &workout.Date, &workout.Exercise); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		setRows, err := db.Query("SELECT weight, reps FROM sets WHERE workout_id = ?", workout.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer setRows.Close()

		for setRows.Next() {
			var set Set
			if err := setRows.Scan(&set.Weight, &set.Reps); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			workout.Sets = append(workout.Sets, set)
		}

		workouts = append(workouts, workout)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	initLogger()
	logger.Info("Starting server", slog.String("port", "8080"))

	initDB("workouts.db")
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/workouts", createWorkoutHandler).Methods("POST")
	r.HandleFunc("/workouts", getWorkoutsHandler).Methods("GET")

	// Wrap the router with the CORS middleware
	http.Handle("/", enableCORS(r))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
