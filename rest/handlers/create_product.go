package handlers

import (
	"encoding/json"
	"net/http"

	"go-first-project/models"
)

// CreateProduct is exported (Capital 'C') so main.go can register it
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input models.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Basic validation
	if input.Name == "" {
		sendData(w, http.StatusBadRequest, map[string]string{"error": "Product name is required"})
		return
	}
	if input.Price <= 0 {
		sendData(w, http.StatusBadRequest, map[string]string{"error": "Product price must be greater than 0"})
		return
	}

	models.Mu.Lock()
	defer models.Mu.Unlock()

	input.ID = models.NextID
	models.NextID++
	models.Products = append(models.Products, input)

	sendData(w, http.StatusCreated, input)
}
