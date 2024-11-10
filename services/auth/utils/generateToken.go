package utils

import (
	"fmt"
	"errors"
	// "time"
	"github.com/golang-jwt/jwt/v5"
)


func ParseToken(tokenString string) (*jwt.Token, error) {
	secretKey := "samyak"
	fmt.Println(secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
