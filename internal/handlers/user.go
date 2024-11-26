package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"log"
	"RTC/internal/models"
	"RTC/internal/validators"

	"github.com/jackc/pgx/v5"
)

func RegisterHandler(db *pgx.Conn) (http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request){
		var req models.RegistrationRequest

		// Parse the JSON into req to later check if is the right input 
		if err := decodeAndValidateRequest(r, &req); err != nil{
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Save user in database
		err := models.SaveUser(db, req)
		if err != nil{
			// Catch and inform the user if already user with said email, so we use costum error
			if errors.Is(err, models.ErrDupEmail) || errors.Is(err, models.ErrDupUsername){
				respondWithError(w, http.StatusConflict, err.Error())
			} else {
				log.Printf("Error saving user:", err)
				respondWithError(w, http.StatusConflict, err.Error())
			}
			return
		} 

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered  successfully"))
	}
}

func decodeAndValidateRequest(r *http.Request, req *models.RegistrationRequest)(error){
	// Parse the JSON into req to later check if is the right input 
	if err := json.NewDecoder(r.Body).Decode(req); err != nil{
		return errors.New("invalid JSON payload")
	}
	// Validate input 
	validator := validators.NewValidator()
	if err := validator.ValidateStruct(req); err != nil{
		return errors.New("invalid input "+ validators.FormatValidationErrors(err))
	}
	return nil
}

// Writes an error response so that it has a consistent format.
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}