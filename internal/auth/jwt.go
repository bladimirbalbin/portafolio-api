package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret []byte

func GenerateToken(user string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(JwtSecret)
}

func ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return JwtSecret, nil
	})
}
