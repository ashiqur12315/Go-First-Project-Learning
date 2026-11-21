package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go-first-project/models"
)

// DeleteProduct is exported (Capital 'D') so main.go can register it
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	models.Mu.Lock()
	defer models.Mu.Unlock()

	for i, p := range models.Products {
		if p.ID == id {
			models.Products = append(models.Products[:i], models.Products[i+1:]...)
			sendData(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Product %d deleted successfully", id)})
			return
		}
	}

	sendData(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
}
