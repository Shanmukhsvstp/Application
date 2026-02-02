package tools

import (
	"errors"
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
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
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

func GetDataFromToken(tokenStr string) (string, float64, float64, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", 0, 0, err // expired or invalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", 0, 0, errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", 0, 0, errors.New("invalid sub claim")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", 0, 0, errors.New("invalid exp claim")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return "", 0, 0, errors.New("invalid iat claim")
	}

	return userID, exp, iat, nil
}
