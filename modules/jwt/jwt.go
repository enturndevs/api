package jwt

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims jwt claims
type Claims map[string]interface{}

// Encode claims to jwt token
func Encode(values Claims, key string) (string, error) {
	claims := jwt.MapClaims{}
	for key, value := range values {
		claims[key] = value
	}

	tokenString, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Decode jwt token to claims
func Decode(tokenString string, key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
