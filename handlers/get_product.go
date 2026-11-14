package handlers

import (
	"net/http"
	"strconv"

	"go-first-project/models"
)

// GetAllProducts is exported (Capital 'G') so main.go can register it
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	models.Mu.RLock()
	defer models.Mu.RUnlock()

	writeJSON(w, http.StatusOK, models.Products)
}

// GetProductByID is exported (Capital 'G') so main.go can register it
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	models.Mu.RLock()
	defer models.Mu.RUnlock()

	for _, p := range models.Products {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}
