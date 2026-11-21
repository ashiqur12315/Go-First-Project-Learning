package config

import (
	"bufio"
	"os"
	"strings"
)

// Config holds all application configuration values.
type Config struct {
	Version     string
	ServiceName string
	Port        string
}

// LoadEnv reads a .env file and sets each key-value pair into the environment using os.Setenv.
// It silently skips if the file does not exist (useful for production or Docker environments).
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No .env file found; proceed with system environment variables
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comment lines starting with '#'
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on the first '=' character
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Strip optional wrapping quotes (e.g. "value" or 'value')
			value = strings.Trim(value, `"'`)

			// Only set into environment if not already set by the operating system
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}

	return scanner.Err()
}

// getEnv retrieves an environment variable or returns a fallback default if not set.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// Load loads the .env file (if present) and returns a populated Config struct with fallback defaults.
func Load() *Config {
	_ = LoadEnv(".env")

	return &Config{
		Version:     getEnv("VERSION", "1.0.0"),
		ServiceName: getEnv("SERVICE_NAME", "product-service"),
		Port:        getEnv("PORT", "8080"),
	}
}
