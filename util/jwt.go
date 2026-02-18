package util

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func SignJwt(secret string, claims jwt.MapClaims) (*string, *AppError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return nil, Error(err, 500, "Failed to sign token")
	}

	return &tokenString, nil
}

func ParseJwt(secret string, tokenString string) (jwt.MapClaims, *AppError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, Error(err, http.StatusInternalServerError, "Failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, Error(err, http.StatusForbidden, "Invalid claims")
	}

	return claims, nil
}
