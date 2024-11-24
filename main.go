package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

type RegistrationRequest struct{
	Email		string `json:"email" validate:"required,email"`
	Password 	string `json:"password" validate:"required,min=8"`
}

var validate = validator.New()
var ErrDuoEmail = errors.New("email already in use")

func main() {
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .en file: %v", err)
	}

	// Constructing the connection string
	dsn := fmt.Sprintf("postgress://%s:%s@%s:%s/%s?sslmode=disable", 
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_PORT"), 
		os.Getenv("DB_NAME"))

	// Connect to the database
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database")

	defer db.Close(context.Background())

	// Set up the router
	r := chi.NewRouter()
	r.Post("/register", registerHandler)


	fmt.Println("Server running on http://localhost:5500")
	go func() {
		log.Fatal(http.ListenAndServe(":5500", nil))
	}()
}

func registerHandler(w http.ResponseWriter, r *http.Request){
	var req RegistrationRequest

	// Parse the JSON into req to later check if is the right input 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Invalid JSON")
		return
	}

	// Validate input 
	if err := validate.Struct(req); err != nil{
		http.Error(w, "Unable to validate input", http.StatusBadRequest)
		log.Println("Unable to validate input")
		return
	}

	// Hash the user password to save it into the database
	hashedPassword := hashPassword(req.Password)

	// Save user in database
	err := saveUser(req.Email, hashedPassword)
	if err != nil{
		// Catch and inform the user if already user with said email, so we use costum error

	} 



}

func hashPassword(password string)(hashedpassword string){

}

func saveUser(email, hashedPassword string)(err error){

}

