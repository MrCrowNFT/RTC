package main

import (
	"log"
	"net/http"
	"RTC/config"
	"RTC/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load config
	config.initConfig()

	// Database connection
	db := config.initDB()
	defer db.Close()

	// Set up the router
	r := chi.NewRouter()
	r.Post("/register", handlers.registerHandler(db))


	log.Println("Server running on http://localhost:5500")
	// Assign router to server
	log.Fatal(http.ListenAndServe(":5500", r))
}
}

