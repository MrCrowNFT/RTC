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

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *pgx.Conn

type RegistrationRequest struct{
	Username	string `json:"username" validate:"required, min=1"`
	Email		string `json:"email" validate:"required,email"`
	Password 	string `json:"password" validate:"required,min=8"`
}

var validate = validator.New()
var (
	ErrDupEmail = errors.New("email already in use")
	ErrDupUsername = errors.New("username already in use")
)

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
	// Assign router to server
	log.Fatal(http.ListenAndServe(":5500", r))
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
		http.Error(w, "Invalid input: "+ formatValidationErrors(err), http.StatusBadRequest)
		log.Println("Validation failed: ", err)
		return
	}

	// Hash the user password to save it into the database
	hashedPassword := hashPassword(req.Password)

	// Save user in database
	err := saveUser(req.Username, req.Email, hashedPassword)
	if err != nil{
		// Catch and inform the user if already user with said email, so we use costum error
		if errors.Is(err, ErrDupEmail) || errors.Is(err, ErrDupUsername){
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Error saving user:", err)
		}
		return
	} 

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User regustered successfully"))
}

func hashPassword(password string)(hashedpassword string){
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func saveUser(username, email, hashedPassword string)(err error){
	_, err = db.Exec(context.Background(),
		`INSERT INTO accounts (username, email, password) VALUES ($1, $2, $3)`,
		username, email, hashedPassword)

	if err != nil{
		// Check posgress unique constraint violation to catch dups
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr){
			// 23505 is code for unique violation
			if pgErr.Code == "23505"{
				// we check the content of the error to determine which element is dup
				if strings.Contains(pgErr.Detail, "email"){
					return ErrDupEmail
				}
				if strings.Contains(pgErr.Detail, "username"){
					return ErrDupUsername
				}
			}
		}
		return err
	}

	return nil
}

func formatValidationErrors(err error)(string){
	var sb strings.Builder
	for _, err :=  range err.(validator.ValidationErrors){
		sb.WriteString(fmt.Sprintf("Field '%s' failed validatin: %s.", err.Field(), err.Tag()))
	}
	return sb.String()
}

