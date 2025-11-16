package http

import (
	"encoding/json"
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"net/http"
)

func NewErrorBody(w http.ResponseWriter, ct string, err error, status int) {
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	logger.Logger.Println(err)
}
