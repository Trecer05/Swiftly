package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, jwt.ErrTokenMalformed
	} else {
		return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
	}
}