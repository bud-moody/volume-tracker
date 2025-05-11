package main

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Updated initDB to clear existing data during initialization
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "workouts.db")
	if err != nil {
		log.Fatal(err)
	}

	// Clear existing data
	db.Exec("DROP TABLE IF EXISTS sets")
	db.Exec("DROP TABLE IF EXISTS workouts")

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
		log.Fatal(err)
	}
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

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/workouts", createWorkoutHandler).Methods("POST")
	r.HandleFunc("/workouts", getWorkoutsHandler).Methods("GET")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
