package cmd

import (
	"fmt"
	"net/http"

	"go-first-project/handlers"
	"go-first-project/middleware"
)

// RegisterRoutes registers all endpoints onto the provided ServeMux.
// It accepts both mux and manager as input parameters so routes can be registered
// and route-specific middleware can be applied where needed.
func RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	// Hello test endpoint
	mux.Handle("GET /hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	}))

	// Product REST API endpoints
	mux.Handle("GET /products", http.HandlerFunc(handlers.GetAllProducts))
	mux.Handle("GET /products/{id}", http.HandlerFunc(handlers.GetProductByID))
	mux.Handle("POST /products", http.HandlerFunc(handlers.CreateProduct))
	mux.Handle("PUT /products/{id}", http.HandlerFunc(handlers.UpdateProduct))
	mux.Handle("PATCH /products/{id}", http.HandlerFunc(handlers.PatchProduct))
	mux.Handle("DELETE /products/{id}", http.HandlerFunc(handlers.DeleteProduct))
}
