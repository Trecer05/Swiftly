package cloud

import (
	"net/http"
	"time"

	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	middleware "github.com/Trecer05/Swiftly/internal/handler"

	"github.com/gorilla/mux"
)

func InitCloudRoutes(r *mux.Router, manager *manager.Manager) {
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	apiSecure := r.PathPrefix("/api/v1/cloud").Subrouter()
	apiSecure.Use(middleware.AuthMiddleware())
	apiSecure.Use(middleware.RateLimitMiddleware(rateLimiter))
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id}/file", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/team/{id}/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/file", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/user/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id}/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id}/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/user/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/team/{id}/file/{file_id}/download", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/file/{file_id}/download", func(w http.ResponseWriter, r *http.Request) {
		
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id}/file/{file_id}/name", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/file/{file_id}/name", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id}/file/{file_id}/share", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPatch)

	apiSecure.HandleFunc("/user/file/{file_id}/share", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPatch)

	apiSecure.HandleFunc("/team/{id}/folder", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/user/folder", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/team/{id}/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id}/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/user/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {
		
	}).Methods(http.MethodDelete)
}
