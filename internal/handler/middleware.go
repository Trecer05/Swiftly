package handler

import (
	"log"
	"net/http"
	"strings"

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
		login := claims["authorization"].(bool)

		if !login {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Unauthorized request")
			return
		}

		next.ServeHTTP(w, r)
	})
}