package utils

import(
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string)(hashedpassword string, err error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}
	return string(hash), nil
} 