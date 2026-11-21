package handlers

import (
	"encoding/json"
	"net/http"

	"go-first-project/database"
	"go-first-project/models"
	"go-first-project/utils"
)

// CreateProduct decodes the JSON payload, validates it, and inserts it into the database.
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input models.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Basic validation
	if input.Name == "" {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Product name is required"})
		return
	}
	if input.Price <= 0 {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Product price must be greater than 0"})
		return
	}

	created := database.CreateProduct(input)
	utils.SendData(w, http.StatusCreated, created)
}
