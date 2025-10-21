package auth

import (
	"encoding/json"
	"net/http"
	"time"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	tokenErrors "github.com/Trecer05/Swiftly/internal/errors/tokens"
	middleware "github.com/Trecer05/Swiftly/internal/handler"
	model "github.com/Trecer05/Swiftly/internal/model/auth"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"
	tokens "github.com/Trecer05/Swiftly/internal/service/auth"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func InitAuthRoutes(router *mux.Router, mgr *manager.Manager) {
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	api := router.PathPrefix("/api/v1").Subrouter()

	apiSecure := router.PathPrefix("/api/v1").Subrouter()
	apiSecure.Use(middleware.AuthMiddleware())
	apiSecure.Use(middleware.RateLimitMiddleware(rateLimiter))

	api.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, mgr)
	}).Methods(http.MethodPost)

	api.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		Register(w, r, mgr)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		Logout(w, r, mgr)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		Refresh(w, r, mgr)
	}).Methods(http.MethodPost)
}

func Login(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var user model.User
	var err error

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err := mgr.Login(&user); err != nil {
		switch err {
		case errors.ErrInvalidPassword:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusUnauthorized)
			return
		case errors.ErrNoUser:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		default:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	var refreshToken string
	token, err := tokens.ValidateRefreshToken(mgr, user.ID)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			refreshToken, err = tokens.GenerateRefreshToken(mgr, user.ID)
			if err != nil {
				serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
				return
			}
		} else {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}
	if token.Token == "" {
		token.Token = refreshToken
	}

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": tokens.AddAccessTime(),
	}

	var accessToken string
	if accessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": token.Token,
	})
}

func Register(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var user model.User
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err = mgr.Register(&user); err != nil {
		switch err {
		case errors.ErrEmailExists:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusConflict)
			return
		case errors.ErrNumberExists:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusConflict)
			return
		default:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	refreshToken, err := tokens.GenerateRefreshToken(mgr, user.ID)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": tokens.AddAccessTime(),
	}

	var accessToken string
	if accessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Logout(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var err error
	var token model.RefreshToken
	var id int
	var ok bool

	if id, ok = r.Context().Value("id").(int); !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	if token, err = tokens.ValidateRefreshToken(mgr, id); err != nil {
		if err == jwt.ErrTokenExpired {
			serviceHttp.NewErrorBody(w, "application/json", tokenErrors.ErrRefreshTokenExpired, http.StatusUnauthorized)
			return
		}
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if err := mgr.DeleteRefreshToken(id, token.Token); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"logout": "success",
	})
}

func Refresh(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var t model.Tokens
	var id int
	var ok bool

	if id, ok = r.Context().Value("id").(int); !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	refreshToken, err := tokens.ValidateRefreshToken(mgr, id)
	if err != nil {
		switch err {
		case jwt.ErrTokenExpired:
			serviceHttp.NewErrorBody(w, "application/json", tokenErrors.ErrRefreshTokenExpired, http.StatusUnauthorized)
			return
		default:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}
	t.RefreshToken = refreshToken.Token

	claims := jwt.MapClaims{
		"id":  id,
		"exp": tokens.AddAccessTime(),
	}

	if t.AccessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
	})
}
