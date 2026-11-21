package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendData(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{"message": "success"}

	SendData(rr, http.StatusOK, payload)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type 'application/json', got %q", contentType)
	}

	var result map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["message"] != "success" {
		t.Fatalf("expected message 'success', got %q", result["message"])
	}
}
