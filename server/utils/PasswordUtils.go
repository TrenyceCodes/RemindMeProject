package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal("There is a problem hashing password: ", err)
		return string(hashPassword), "There is a problem hashing password: ", err
	}

	return string(hashPassword), "Password hashed successfully", err
}

func ValidatePassword(password string, userPassword string) (bool, string) {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		return false, "Invalid password"
	}

	return true, "Password is valid"
}
