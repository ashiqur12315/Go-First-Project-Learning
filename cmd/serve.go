package cmd

import (
	"fmt"
	"net/http"

	"go-first-project/middleware"
)

// Serve configures the router, attaches middleware, and starts the HTTP server.
// Exported (Capital 'S') so it can be called from main.go.
func Serve() {
	mux := http.NewServeMux()
	manager := middleware.NewManager()

	// 1. Configure global middlewares
	// Execution flow: Incoming Request -> LoggerMiddleware -> CorsMiddleware -> Mux -> Handlers
	manager.Use(
		middleware.LoggerMiddleware,
		middleware.CorsMiddleware,
	)

	// 2. Register all application routes, passing both mux and manager as input
	RegisterRoutes(mux, manager)

	// 3. Wrap the router with all registered middlewares
	rootHandler := manager.Wrap(mux)

	// 4. Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", rootHandler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
