package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseRecorder is a custom wrapper around http.ResponseWriter.
// It intercepts calls to WriteHeader so we can capture the HTTP status code.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before passing it to the real ResponseWriter
func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// LoggerMiddleware logs every incoming HTTP request with:
// - HTTP Method (GET, POST, etc.)
// - Request Path (/products, /hello, etc.)
// - Response Status Code (200, 201, 404, etc.)
// - Execution Duration (latency)
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap 'w' with responseRecorder (defaults to 200 OK)
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Pass the request to the next middleware or handler in the chain
		next.ServeHTTP(rec, r)

		// Calculate how long the entire request took
		duration := time.Since(start)

		// Print formatted log output to the terminal
		log.Printf("[%s] %s | %d | %s", r.Method, r.URL.Path, rec.statusCode, duration)
	})
}
