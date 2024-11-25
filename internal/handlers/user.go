package handlers

import (
	"encoding/json"
	"net/http"
	"RTC/internal/models"
	"RTC/internal/validators"

	"github.com/jackc/pgx/v5"
)

func registerHandler(db *pgx.Conn) (http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request){
		var req models.RegistrationRequest

		// Parse the JSON into req to later check if is the right input 
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate input 
		if err := validators.ValidateStruct(req); err != nil{
			http.Error(w, "Invalid input: "+ validators.formatValidationErrors(err), http.StatusBadRequest)
			return
		}

		// Hash the user password to save it into the database
		hashedPassword, err := hashPassword(req.Password)
		if err != nil{
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
    		log.Printf("Error hashing password:", err)
    		return
		}

		// Save user in database
		err = saveUser(req.Username, req.Email, hashedPassword)
		if err != nil{
			// Catch and inform the user if already user with said email, so we use costum error
			if errors.Is(err, ErrDupEmail) || errors.Is(err, ErrDupUsername){
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Printf("Error saving user:", err)
			}
			return
		} 

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User regustered successfully"))
	}
}