package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(port string, r *mux.Router) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}