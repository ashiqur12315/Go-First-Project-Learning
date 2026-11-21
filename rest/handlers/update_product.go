package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-first-project/models"
)

// UpdateProduct is exported (Capital 'U') so main.go can register it
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	var input models.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Basic validation for full replacement
	if input.Name == "" || input.Price <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Product name and positive price are required"})
		return
	}

	models.Mu.Lock()
	defer models.Mu.Unlock()

	for i, p := range models.Products {
		if p.ID == id {
			input.ID = id
			models.Products[i] = input
			writeJSON(w, http.StatusOK, models.Products[i])
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}

// PatchProduct is exported (Capital 'P') so main.go can register it
func PatchProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	var input models.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	models.Mu.Lock()
	defer models.Mu.Unlock()

	for i, p := range models.Products {
		if p.ID == id {
			if input.Name != nil {
				models.Products[i].Name = *input.Name
			}
			if input.Description != nil {
				models.Products[i].Description = *input.Description
			}
			if input.Price != nil {
				if *input.Price <= 0 {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be greater than 0"})
					return
				}
				models.Products[i].Price = *input.Price
			}
			if input.Stock != nil {
				if *input.Stock < 0 {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Stock cannot be negative"})
					return
				}
				models.Products[i].Stock = *input.Stock
			}

			writeJSON(w, http.StatusOK, models.Products[i])
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}
