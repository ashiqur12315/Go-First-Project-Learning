package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go-first-project/database"
	"go-first-project/models"
	"go-first-project/utils"
)

// UpdateProduct handles full replacement of a product (PUT).
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	var input models.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Basic validation for full replacement
	if input.Name == "" || input.Price <= 0 {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Product name and positive price are required"})
		return
	}

	updated, err := database.UpdateProduct(id, input)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			utils.SendData(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
			return
		}
		utils.SendData(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
		return
	}

	utils.SendData(w, http.StatusOK, updated)
}

// PatchProduct handles partial updates to a product (PATCH).
func PatchProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	var input models.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Validation on provided partial fields
	if input.Price != nil && *input.Price <= 0 {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Price must be greater than 0"})
		return
	}
	if input.Stock != nil && *input.Stock < 0 {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Stock cannot be negative"})
		return
	}

	patched, err := database.PatchProduct(id, input)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			utils.SendData(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
			return
		}
		utils.SendData(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
		return
	}

	utils.SendData(w, http.StatusOK, patched)
}
