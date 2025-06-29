package auth

import (
	"encoding/json"
	"log"
	"net/http"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	model "github.com/Trecer05/Swiftly/internal/model/auth"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"
	tokens "github.com/Trecer05/Swiftly/internal/service/auth"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gorilla/mux"
)

func InitAuthRoutes(router *mux.Router, mgr *manager.Manager) {
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, mgr)
	}).Methods(http.MethodPost)

	api.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		Register(w, r, mgr)
	}).Methods(http.MethodPost)
}

func Login(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var user model.User
	var err error

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := mgr.Login(&user); err != nil {
		switch err {
		case errors.ErrInvalidPassword:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			log.Println(err)
			return
		case errors.ErrNoUser:
			http.Error(w, err.Error(), http.StatusNotFound)
			log.Println(err)
			return
		}
	}

	var refreshToken string
	if err := tokens.ValidateRefreshToken(mgr, user.ID); err != nil {
		if err == jwt.ErrTokenExpired {
			refreshToken, err = tokens.GenerateRefreshToken(mgr, user.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     tokens.AddAccessTime(),
	}

	var accessToken string
	if accessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Register(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var user model.User
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err = mgr.Register(&user); err != nil {
		if err == errors.ErrEmailExists {
			http.Error(w, err.Error(), http.StatusConflict)
			log.Println(err)
			return
		} else if err == errors.ErrNumberExists {
			http.Error(w, err.Error(), http.StatusConflict)
			log.Println(err)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	refreshToken, err := tokens.GenerateRefreshToken(mgr, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     tokens.AddAccessTime(),
	}

	var accessToken string
	if accessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}