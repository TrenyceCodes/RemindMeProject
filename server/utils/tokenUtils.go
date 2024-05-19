package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func CreateToken(username, userId string) (string, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Failed to load environment variables")
		return "", err
	}

	secretKey := envFile["SECRET_KEY"]
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set in the environment variables")
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
