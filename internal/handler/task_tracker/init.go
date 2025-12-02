package task_tracker

import (
	"net/http"
	"time"

	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	middleware "github.com/Trecer05/Swiftly/internal/handler"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"

	"github.com/gorilla/mux"
)

func InitTaskRoutes(r *mux.Router, manager *manager.Manager, rds *redis.WebSocketManager) {
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	apiSecure := r.PathPrefix("/api/v1").Subrouter()
	apiSecure.Use(middleware.AuthMiddleware())
	apiSecure.Use(middleware.RateLimitMiddleware(rateLimiter))
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	
	apiSecure.HandleFunc("/team/{id}/dashboard", func(w http.ResponseWriter, r *http.Request) {
		DashboardWSHandler(w, r, manager, rds)
	})
	
	apiSecure.HandleFunc("/team/{id}/dashboard/info", func(w http.ResponseWriter, r *http.Request) {
		GetDashboardInfoHandler(w, r, manager)
	}).Methods(http.MethodGet)
	
	apiSecure.HandleFunc("/team/{team_id}/column/{column_id}/task", func(w http.ResponseWriter, r *http.Request) {
		CreateTaskHandler(w, r, manager, rds)
	}).Methods(http.MethodPost)
	
	apiSecure.HandleFunc("/team/{team_id}/column/{column_id}/task/{task_id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteTaskHandler(w, r, manager, rds)
	}).Methods(http.MethodDelete)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}", func(w http.ResponseWriter, r *http.Request) {
		GetTaskHandler(w, r, manager)
	}).Methods(http.MethodGet)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}/title", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskTitleHandler(w, r, manager, rds)
	}).Methods(http.MethodPut)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}/description", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskDescriptionHandler(w, r, manager, rds)
	}).Methods(http.MethodPut)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}/deadline", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskDeadlineHandler(w, r, manager, rds)
	}).Methods(http.MethodPut)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}/developer", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskDeveloperHandler(w, r, manager, rds)
	}).Methods(http.MethodPut)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}/status", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskStatusHandler(w, r, manager, rds)
	}).Methods(http.MethodPut)
	
	apiSecure.HandleFunc("/team/{team_id}/task/{task_id}/move", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskColumnHandler(w, r, manager, rds)
	}).Methods(http.MethodPut)
}
