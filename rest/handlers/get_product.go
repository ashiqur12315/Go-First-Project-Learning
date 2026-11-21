package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"go-first-project/database"
	"go-first-project/utils"
)

// GetAllProducts fetches all products from the database.
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products := database.GetAllProducts()
	utils.SendData(w, http.StatusOK, products)
}

// GetProductByID retrieves a single product by its ID from the database.
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	product, err := database.GetProductByID(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			utils.SendData(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
			return
		}
		utils.SendData(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch product"})
		return
	}

	utils.SendData(w, http.StatusOK, product)
}
