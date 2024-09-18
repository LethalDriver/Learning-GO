package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func parseRequest(r *http.Request, reqStruct any) error {
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        return fmt.Errorf("failed reading body: %w", err)
    }

    err = json.Unmarshal(bodyBytes, reqStruct) // Unmarshal into the pointer
    if err != nil {
        return fmt.Errorf("failed parsing body: %w", err)
    }
	
    return nil
}

func writeResponse(w http.ResponseWriter, respStruct any) error {
	tokenJson, err := json.Marshal(respStruct)
	if err != nil {
		return fmt.Errorf("failed marshaling response: %v", err)
	}
	_, err = w.Write(tokenJson); 
	if err != nil {
		return fmt.Errorf("failed writing to response: %v", err)
	}
	return nil
}