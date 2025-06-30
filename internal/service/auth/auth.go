package auth

import (
	"database/sql"
	"os"
	"time"

	models "github.com/Trecer05/Swiftly/internal/model/auth"
	"github.com/Trecer05/Swiftly/internal/repository/postgres/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = os.Getenv("JWT_SECRET")

func ParseToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, jwt.ErrTokenMalformed
	} else {
		return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	}
}

func GenerateRefreshToken(manager *auth.Manager, userId int) (string, error) {
	token := uuid.New().String()

	err := manager.SaveRefreshToken(token, userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateAccessToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AddAccessTime() int64 {
	return time.Now().Add(15 * time.Minute).Unix()
}

func ValidateRefreshToken(manager *auth.Manager, userId int) (models.RefreshToken, error) {
	token, err := manager.GetRefreshToken(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.RefreshToken{}, jwt.ErrTokenExpired
		}
		return token, err
	}

	if time.Now().After(token.ExpiredAt) {
		if err := manager.DeleteRefreshToken(userId, token.Token); err != nil {
			return models.RefreshToken{}, err
		} else {
			return models.RefreshToken{}, jwt.ErrTokenExpired
		}
	}
	return token, nil
}
