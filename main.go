package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// Product represents our data model for a product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// UpdateProductInput is used for PATCH requests where fields are optional.
// Pointer types (*string, *float64, *int) allow us to distinguish between
// a field being omitted (nil) versus sent with a zero value ("", 0).
type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

// In-memory product storage with some initial seed data
var (
	mu       sync.RWMutex
	products = []Product{
		{
			ID:          1,
			Name:        "Mechanical Keyboard",
			Description: "RGB mechanical gaming keyboard with blue switches",
			Price:       79.99,
			Stock:       25,
		},
		{
			ID:          2,
			Name:        "Wireless Mouse",
			Description: "Ergonomic 2.4GHz wireless mouse with adjustable DPI",
			Price:       29.99,
			Stock:       50,
		},
		{
			ID:          3,
			Name:        "Gaming Monitor",
			Description: "27-inch 144Hz IPS display with 1ms response time",
			Price:       199.99,
			Stock:       15,
		},
	}
	nextID = 4
)

// Helper function to send standard JSON responses
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// GET /products - Get all products
func getAllProducts(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	writeJSON(w, http.StatusOK, products)
}

// GET /products/{id} - Get a single product by ID
func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	for _, p := range products {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}

// POST /products - Create a new product
func createProduct(w http.ResponseWriter, r *http.Request) {
	var input Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Basic validation
	if input.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Product name is required"})
		return
	}
	if input.Price <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Product price must be greater than 0"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	input.ID = nextID
	nextID++
	products = append(products, input)

	writeJSON(w, http.StatusCreated, input)
}

// PUT /products/{id} - Completely replace/update an existing product
func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	var input Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Basic validation for full replacement
	if input.Name == "" || input.Price <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Product name and positive price are required"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, p := range products {
		if p.ID == id {
			// Replace with new data while preserving the original ID
			input.ID = id
			products[i] = input
			writeJSON(w, http.StatusOK, products[i])
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}

// PATCH /products/{id} - Partially update fields of an existing product
func patchProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	var input UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, p := range products {
		if p.ID == id {
			// Update only fields that were provided in the request body
			if input.Name != nil {
				products[i].Name = *input.Name
			}
			if input.Description != nil {
				products[i].Description = *input.Description
			}
			if input.Price != nil {
				if *input.Price <= 0 {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be greater than 0"})
					return
				}
				products[i].Price = *input.Price
			}
			if input.Stock != nil {
				if *input.Stock < 0 {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Stock cannot be negative"})
					return
				}
				products[i].Stock = *input.Stock
			}

			writeJSON(w, http.StatusOK, products[i])
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}

// DELETE /products/{id} - Delete a product by ID
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, p := range products {
		if p.ID == id {
			// Remove element from slice
			products = append(products[:i], products[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Product %d deleted successfully", id)})
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}

// Example 1: An endpoint where CORS logic is written DIRECTLY inside the handler (NO middleware).
// Look how simple it is: it's just setting headers on 'w' and checking if it's an OPTIONS request!
func exampleDirectCORS(w http.ResponseWriter, r *http.Request) {
	// --- CORS logic starts here ---
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	// --- CORS logic ends here ---

	// The actual endpoint work:
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "This endpoint handles CORS directly inside itself without any middleware!",
	})
}

// corsMiddleware adds CORS headers to responses.
// Notice the step-by-step flow:
// 1. We create the handler function logic.
// 2. We convert that function into an http.Handler.
// 3. We return that newly created http.Handler at the very end.
func corsMiddleware(next http.Handler) http.Handler {
	// Step 1: Write the handler function that intercepts every incoming request
	middlewareHandler := func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight check: if browser sends OPTIONS, finish immediately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue the chain: call the next handler (the router / ServeMux)
		next.ServeHTTP(w, r)
	}

	// Step 2: Convert our function into an http.Handler using the http.HandlerFunc adapter
	wrappedHandler := http.HandlerFunc(middlewareHandler)

	// Step 3: Return the final wrapped handler at the end of the function!
	return wrappedHandler
}

func main() {
	mux := http.NewServeMux()

	// Existing test endpoint
	mux.Handle("GET /hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	}))

	// Product REST API endpoints using mux.Handle and http.HandlerFunc adapter
	mux.Handle("GET /products", http.HandlerFunc(getAllProducts))
	mux.Handle("GET /products/{id}", http.HandlerFunc(getProductByID))
	mux.Handle("POST /products", http.HandlerFunc(createProduct))
	mux.Handle("PUT /products/{id}", http.HandlerFunc(updateProduct))
	mux.Handle("PATCH /products/{id}", http.HandlerFunc(patchProduct))
	mux.Handle("DELETE /products/{id}", http.HandlerFunc(deleteProduct))

	// Demonstration Endpoint 1: CORS written directly inside the endpoint handler (no middleware)
	mux.Handle("/example-direct-cors", http.HandlerFunc(exampleDirectCORS))

	// Demonstration Endpoint 2: Applying the middleware to a SINGLE endpoint on ONE line
	mux.Handle("GET /example-inline-cors", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "This endpoint was wrapped with corsMiddleware directly in one line!",
		})
	})))

	// Wrap our router (mux) with the CORS middleware
	// Now EVERY incoming request goes through corsMiddleware first!
	handlerWithCORS := corsMiddleware(mux)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
