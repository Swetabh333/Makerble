package helper

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateToken(username string, ID string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = username
	claims["sub"] = ID
	claims["role"] = role
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(username string, ID string, role string) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = username
	claims["sub"] = ID
	claims["role"] = role
	claims["exp"] = time.Now().Add(30 * 24 * time.Hour).Unix()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
