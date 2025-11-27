package task_tracker

import (
	"net/http"
	"time"

	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	middleware "github.com/Trecer05/Swiftly/internal/handler"

	"github.com/gorilla/mux"
)

func InitTaskRoutes(r *mux.Router, manager *manager.Manager) {
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	apiSecure := r.PathPrefix("/api/v1").Subrouter()
	apiSecure.Use(middleware.AuthMiddleware())
	apiSecure.Use(middleware.RateLimitMiddleware(rateLimiter))
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	
	// вроде эти два по вебсокетам будут, а гет просто
	apiSecure.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		CreateTaskHandler(w, r, manager)
	}).Methods(http.MethodPost)
	
	apiSecure.HandleFunc("/task/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteTaskHandler(w, r, manager)
	}).Methods(http.MethodDelete)
	
	apiSecure.HandleFunc("/task/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetTaskHandler(w, r, manager)
	}).Methods(http.MethodGet)
}
