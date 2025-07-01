package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	tokenErrors "github.com/Trecer05/Swiftly/internal/errors/tokens"
	auth "github.com/Trecer05/Swiftly/internal/service/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": errors.ErrUnauthorized.Error()})
				log.Println("Unauthorized request")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid authorization header format"})
				log.Println("Invalid authorization header format")
				return
			}
			tokenString := parts[1]

			token, err := auth.ParseToken(tokenString)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid token: " + err.Error()})
				log.Println("Invalid token:", err)
				return
			}
			if !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "token is not valid"})
				log.Println("Token is not valid")
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": tokenErrors.ErrAccessTokenExpired.Error()})
					log.Println("Token expired")
					return
				}
			}

			idFloat, ok := claims["id"].(float64)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid token payload"})
				log.Println("Invalid token payload")
				return
			}
			id := int(idFloat)

			ctx := context.WithValue(r.Context(), "id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}