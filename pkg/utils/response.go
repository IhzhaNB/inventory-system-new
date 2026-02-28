package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Response represents the standard JSON format for all API responses.
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`   // omitempty: omitted if empty
	Meta    any    `json:"meta,omitempty"`   // Used for pagination metadata
	Errors  any    `json:"errors,omitempty"` // Used for specific error details/validation
}

// WriteJSON is a helper function to encode and send JSON responses.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Success sends a standard success response (e.g., 200 OK or 201 Created).
func Success(w http.ResponseWriter, r *http.Request, status int, message string, data any) {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
		Meta: map[string]string{
			"request_id": middleware.GetReqID(r.Context()),
		},
	}
	WriteJSON(w, status, res)
}

// Error sends a standard failure response (e.g., 400 Bad Request, 500 Internal Server Error).
func Error(w http.ResponseWriter, r *http.Request, status int, message string, errors any) {
	res := Response{
		Status:  status,
		Message: message,
		Errors:  errors,
		Meta: map[string]string{
			"request_id": middleware.GetReqID(r.Context()),
		},
	}
	WriteJSON(w, status, res)
}
