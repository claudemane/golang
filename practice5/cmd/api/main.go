package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"practice5/internal/handler"
	"practice5/internal/repository"
	"time"

	_ "github.com/lib/pq"
)

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

	var db *sql.DB
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

	repo := repository.New(db)
	h := handler.New(repo)

	http.HandleFunc("/users", h.GetUsers)
	http.HandleFunc("/users/common-friends", h.GetCommonFriends)

	log.Println("Starting the Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
