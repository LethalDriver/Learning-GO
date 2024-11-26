package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// parseRequest reads the body of an HTTP request and parses it into a struct.
func parseRequest(r *http.Request, reqStruct any) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed reading body: %w", err)
	}

	err = json.Unmarshal(bodyBytes, reqStruct)
	if err != nil {
		return fmt.Errorf("failed parsing body: %w", err)
	}

	return nil
}

// writeJsonResponse writes a JSON response to an HTTP response writer and sets the content type.
func writeJsonResponse(w http.ResponseWriter, respStruct any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonBytes, err := json.Marshal(respStruct)
	if err != nil {
		return fmt.Errorf("failed marshaling response: %v", err)
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed writing to response: %v", err)
	}
	return nil
}
