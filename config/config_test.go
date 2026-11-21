package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := Load()

	if cfg == nil {
		t.Fatal("expected config to not be nil")
	}

	if cfg.Port == "" {
		t.Errorf("expected Port to be set, got empty string")
	}

	if cfg.ServiceName == "" {
		t.Errorf("expected ServiceName to be set, got empty string")
	}

	if cfg.Version == "" {
		t.Errorf("expected Version to be set, got empty string")
	}
}

func TestLoadEnvCustomFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `
# Sample comment
TEST_CUSTOM_KEY="custom_value"
TEST_PORT='9090'
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	if err := LoadEnv(tmpFile.Name()); err != nil {
		t.Fatalf("LoadEnv failed: %v", err)
	}

	if val := os.Getenv("TEST_CUSTOM_KEY"); val != "custom_value" {
		t.Errorf("expected 'custom_value', got %q", val)
	}

	if val := os.Getenv("TEST_PORT"); val != "9090" {
		t.Errorf("expected '9090', got %q", val)
	}
}

func TestGetEnvFallback(t *testing.T) {
	fallback := "fallback_val"
	val := getEnv("NON_EXISTENT_VAR_12345", fallback)
	if val != fallback {
		t.Errorf("expected %q, got %q", fallback, val)
	}
}
