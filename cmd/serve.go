package cmd

import (
	"fmt"
	"net/http"

	"go-first-project/handlers"
	"go-first-project/middleware"
)

// Serve configures the router, attaches middleware, and starts the HTTP server.
// Exported (Capital 'S') so it can be called from main.go.
func Serve() {
	mux := http.NewServeMux()

	// Hello test endpoint
	mux.Handle("GET /hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	}))

	// Product REST API endpoints using the handlers package
	mux.Handle("GET /products", http.HandlerFunc(handlers.GetAllProducts))
	mux.Handle("GET /products/{id}", http.HandlerFunc(handlers.GetProductByID))
	mux.Handle("POST /products", http.HandlerFunc(handlers.CreateProduct))
	mux.Handle("PUT /products/{id}", http.HandlerFunc(handlers.UpdateProduct))
	mux.Handle("PATCH /products/{id}", http.HandlerFunc(handlers.PatchProduct))
	mux.Handle("DELETE /products/{id}", http.HandlerFunc(handlers.DeleteProduct))

	// Chain middlewares:
	// Flow: Incoming Request -> LoggerMiddleware -> CorsMiddleware -> Mux -> Handlers
	rootHandler := middleware.LoggerMiddleware(middleware.CorsMiddleware(mux))

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", rootHandler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
