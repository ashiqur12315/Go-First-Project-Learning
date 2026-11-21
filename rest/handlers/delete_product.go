package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-first-project/database"
	"go-first-project/utils"
)

// DeleteProduct removes a product by its ID from the database.
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendData(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	err = database.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			utils.SendData(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
			return
		}
		utils.SendData(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
		return
	}

	utils.SendData(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Product %d deleted successfully", id)})
}
