package utils

import(
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordMismatch = errors.New("incorrect password")

func HashPassword(password string)(hashedpassword string, err error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}
	return string(hash), nil
} 

func ComparePassword(hashedPassword, plainPassword string)(error){
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil{
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			return ErrPasswordMismatch
		}
		return err
	}
	return err
}