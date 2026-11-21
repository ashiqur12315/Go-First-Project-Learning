package utils

import (
	"encoding/json"
	"net/http"
)

// SendData is a reusable utility function that sends standard JSON responses
// with the appropriate Content-Type header and HTTP status code.
// Exported (Capital 'S') so any package across the application can reuse it.
func SendData(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
