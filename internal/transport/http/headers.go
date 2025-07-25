package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func NewErrorBody(w http.ResponseWriter, ct string, err error, status int) {
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	log.Println(err)
}
