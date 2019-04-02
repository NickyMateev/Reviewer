package web

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse defines an error response payload
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteResponse writes a payload to the provided writer
func WriteResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
