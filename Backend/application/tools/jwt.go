package tools

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte{}

func init() {
	jwtSecretStr := os.Getenv("JWT_SECRET")

	if jwtSecretStr != "" {
		log.Println("Using JWT secret from environment variable")
	} else {
		jwtSecretStr = "nothing_here"
		log.Println("JWT Not detected in ENV!")
	}

	jwtSecret = []byte(jwtSecretStr)
}

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"expires": time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
