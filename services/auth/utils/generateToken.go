package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = "Samyak"

func ParseToken(tokenString string) (*jwt.Token, error) {
	secretKey := JWT_SECRET

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


func GenerateVerificationToken(voterID string) (string, error) {
	// Create token claims with the voter ID and expiration time
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": voterID,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	// Sign the token with the secret key
	token, err := claims.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err // Return the error to be handled by the calling function
	}

	return token, nil
}
