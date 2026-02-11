package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET")

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token *string) (int64, error) {
	parsedToken, err := jwt.Parse(*token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("Could not parse token")
	}

	isValidToken := parsedToken.Valid
	if !isValidToken {
		return 0, errors.New("Invalid Token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	userId := int64(claims["userId"].(float64))
	return userId, nil
}
