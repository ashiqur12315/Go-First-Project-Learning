package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /products", http.HandlerFunc(getAllProducts))
	mux.Handle("GET /products/{id}", http.HandlerFunc(getProductByID))
	mux.Handle("POST /products", http.HandlerFunc(createProduct))
	mux.Handle("PUT /products/{id}", http.HandlerFunc(updateProduct))
	mux.Handle("PATCH /products/{id}", http.HandlerFunc(patchProduct))
	mux.Handle("DELETE /products/{id}", http.HandlerFunc(deleteProduct))
	return mux
}

func TestGetAllProducts(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var prods []Product
	if err := json.Unmarshal(rr.Body.Bytes(), &prods); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(prods) == 0 {
		t.Fatalf("expected products to be non-empty")
	}
}

func TestGetProductByID(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var prod Product
	if err := json.Unmarshal(rr.Body.Bytes(), &prod); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if prod.ID != 1 {
		t.Fatalf("expected ID 1, got %d", prod.ID)
	}
}

func TestCreateProduct(t *testing.T) {
	server := setupTestServer()
	body := `{"name": "Gaming Headset", "description": "7.1 surround sound", "price": 49.99, "stock": 20}`
	req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var prod Product
	if err := json.Unmarshal(rr.Body.Bytes(), &prod); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if prod.Name != "Gaming Headset" || prod.ID == 0 {
		t.Fatalf("unexpected created product: %+v", prod)
	}
}

func TestUpdateProduct(t *testing.T) {
	server := setupTestServer()
	body := `{"name": "Updated Keyboard", "description": "New description", "price": 99.99, "stock": 30}`
	req := httptest.NewRequest("PUT", "/products/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var prod Product
	if err := json.Unmarshal(rr.Body.Bytes(), &prod); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if prod.Name != "Updated Keyboard" || prod.Price != 99.99 {
		t.Fatalf("unexpected updated product: %+v", prod)
	}
}

func TestPatchProduct(t *testing.T) {
	server := setupTestServer()
	// Only update price
	body := `{"price": 109.99}`
	req := httptest.NewRequest("PATCH", "/products/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var prod Product
	if err := json.Unmarshal(rr.Body.Bytes(), &prod); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if prod.Price != 109.99 || prod.Name != "Updated Keyboard" {
		t.Fatalf("unexpected patched product: %+v", prod)
	}
}

func TestDeleteProduct(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest("DELETE", "/products/2", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	// Verify it cannot be fetched
	req2 := httptest.NewRequest("GET", "/products/2", nil)
	rr2 := httptest.NewRecorder()
	server.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr2.Code)
	}
}
