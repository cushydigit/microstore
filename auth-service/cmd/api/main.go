package main

import (
	"log"
	"net/http"

	"github.com/cushydigit/microstore/auth-service/internal/database"
	"github.com/cushydigit/microstore/auth-service/internal/handlers"
)


func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/register", handlers.Register(db))
	mux.HandleFunc("/login", handlers.Login(db))

	log.Println("Starting Auth Service on :8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
