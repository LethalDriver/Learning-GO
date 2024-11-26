package handler

import (
	"encoding/json"
	"net/http"
)

// parseRequest reads the body of an HTTP request and parses it into a struct.
func writeJsonResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
