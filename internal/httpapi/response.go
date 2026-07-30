package httpapi

import (
	"encoding/json"
	"net/http"
	"time"
)

type Meta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

type Envelope struct {
	Data  any    `json:"data"`
	Meta  Meta   `json:"meta"`
	Error any    `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Data: data,
		Meta: Meta{
			RequestID: "req_local", // fine for now; real IDs later
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Error: nil,
	})
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Data: nil,
		Meta: Meta{
			RequestID: "req_local",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Error: APIError{Code: code, Message: message, Details: []any{}},
	})
}