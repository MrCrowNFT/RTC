package models

import (
	"context"
	"errors"
	"RTC/pkg/hashing"

	"github.com/jackc/pgx/v5"
)

type RegistrationRequest struct{
	Username	string `json:"username" validate:"required,min=1"`
	Email		string `json:"email" validate:"required,email"`
	Password 	string `json:"password" validate:"required,min=8"`
}

var (
	ErrDupEmail = errors.New("email already in use")
	ErrDupUsername = errors.New("username already in use")
)

func saveUser(db *pgx.Conn, req RegistrationRequest)(err error){
	hashedPassword := hashing.HashPassword(req.Password)
	
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
