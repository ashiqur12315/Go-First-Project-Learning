package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON is unexported (starts with lowercase 'w') because it is only
// needed internally by files within the 'handlers' package.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
