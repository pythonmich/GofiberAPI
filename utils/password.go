package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword will hash our users password before its stored in a database
func HashPassword(password string) (string,error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); if err != nil{
		return "", errors.New("unable to hash password")
	}
	return string(hashPassword), nil
}

// CheckPassword will compare passwords to check if the match will return error if no match
func CheckPassword(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}