package rest

import (
	"fmt"
	"net/http"

	"go-first-project/config"
	"go-first-project/rest/middleware"
)

// Start initializes the router, attaches middleware, and starts the HTTP server
// using the configuration passed from cmd.
func Start(cfg *config.Config) {
	mux := http.NewServeMux()
	manager := middleware.NewManager()

	// 1. Configure global middlewares
	// Execution flow: Incoming Request -> LoggerMiddleware -> CorsMiddleware -> Mux -> Handlers
	manager.Use(
		middleware.LoggerMiddleware,
		middleware.CorsMiddleware,
	)

	// 2. Register all REST API endpoints
	RegisterRoutes(mux, manager)

	// 3. Wrap the router with all registered middlewares
	rootHandler := manager.Wrap(mux)

	// 4. Start the HTTP server on configured port
	addr := ":" + cfg.Port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, rootHandler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
