package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Trecer05/Swiftly/internal/service/auth"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Unauthorized request")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusBadRequest)
			log.Println("Invalid authorization header format")
			return
		}
		tokenString := parts[1]

		token, err := auth.ParseToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Println("Invalid token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				log.Println("Token expired")
				return
			}
		}

		idFloat, ok := claims["id"].(float64)
		if !ok {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}
		id := int(idFloat)

		ctx := context.WithValue(r.Context(), "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}