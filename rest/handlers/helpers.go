package handlers

import (
	"net/http"

	"go-first-project/utils"
)

// sendData is a handler-level utility function that reuses utils.SendData.
func sendData(w http.ResponseWriter, status int, data any) {
	utils.SendData(w, status, data)
}
