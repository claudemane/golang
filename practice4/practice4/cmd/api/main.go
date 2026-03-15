package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name,
	)

	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Println("Waiting for database...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	log.Println("Database connected")
	log.Println("Starting the Server on :8080")

	http.HandleFunc("/tasks", tasksHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		rows, err := db.Query("SELECT id, title, done FROM tasks")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var t Task
			rows.Scan(&t.ID, &t.Title, &t.Done)
			tasks = append(tasks, t)
		}

		if tasks == nil {
			tasks = []Task{}
		}
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var t Task
		json.NewDecoder(r.Body).Decode(&t)

		err := db.QueryRow(
			"INSERT INTO tasks(title, done) VALUES($1, false) RETURNING id",
			t.Title,
		).Scan(&t.ID)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)

	case http.MethodPut:
		var t Task
		json.NewDecoder(r.Body).Decode(&t)

		_, err := db.Exec(
			"UPDATE tasks SET title=$1, done=$2 WHERE id=$3",
			t.Title, t.Done, t.ID,
		)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"updated": true})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")

		_, err := db.Exec("DELETE FROM tasks WHERE id=$1", id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"deleted": true})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
