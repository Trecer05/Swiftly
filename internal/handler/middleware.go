package handler

import (
	"context"
	"encoding/json"
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"net/http"
	"strings"
	"sync"
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
				logger.Logger.Println("Unauthorized request")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid authorization header format"})
				logger.Logger.Println("Invalid authorization header format")
				return
			}
			tokenString := parts[1]

			token, err := auth.ParseToken(tokenString)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": tokenErrors.ErrAccessTokenExpired.Error()})
				logger.Logger.Println("Invalid token:", err)
				return
			}
			if !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "token is not valid"})
				logger.Logger.Println("Token is not valid")
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": tokenErrors.ErrAccessTokenExpired.Error()})
					logger.Logger.Println("Token expired")
					return
				}
			}

			idFloat, ok := claims["id"].(float64)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid token payload"})
				logger.Logger.Println("Invalid token payload")
				return
			}
			id := int(idFloat)

			ctx := context.WithValue(r.Context(), "id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *RateLimiter) Allow(identifier string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    now := time.Now()
    cutoff := now.Add(-rl.window)
    
    var validRequests []time.Time
    for _, reqTime := range rl.requests[identifier] {
        if reqTime.After(cutoff) {
            validRequests = append(validRequests, reqTime)
        }
    }
    
    if len(validRequests) >= rl.limit {
        return false
    }
    
    validRequests = append(validRequests, now)
    rl.requests[identifier] = validRequests
    
    return true
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
                ip = forwarded
            }
            
            if !rl.Allow(ip) {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
