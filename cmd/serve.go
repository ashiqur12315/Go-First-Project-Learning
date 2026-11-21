package cmd

import (
	"fmt"

	"go-first-project/config"
	"go-first-project/rest"
)

// Serve loads configuration and starts the REST server.
// Exported (Capital 'S') so it can be called from main.go.
func Serve() {
	// 1. Load configuration from .env file / environment variables
	cfg := config.Load()

	// 2. Print configuration details loaded from .env
	fmt.Println("===========================================")
	fmt.Printf("Service: %s\n", cfg.ServiceName)
	fmt.Printf("Version: %s\n", cfg.Version)
	fmt.Printf("Port:    %s\n", cfg.Port)
	fmt.Println("===========================================")

	// 3. Start the REST server with the loaded configuration
	rest.Start(cfg)
}
