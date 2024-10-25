package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var style string

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, age TEXT)`)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	style = loadAsset("assets/style.css")

	initDB()

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			cleanExpiredSessions()
		}
	}()

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/create", createHandler)

	http.HandleFunc("/read", readHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", pageNotFoundHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
