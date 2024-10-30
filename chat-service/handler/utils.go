package handler

import (
	"encoding/json"
	"net/http"
)

func writeJsonResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
