package http

import (
	"encoding/json"
	"net/http"
)


// errorResponse contains fields of response when server must return error
type errorResponse struct {
	OK          bool              `json:"ok"`
	ErrorCode   int               `json:"error_code,omitempty"`
	Description string            `json:"description,omitempty"`
}

// NewErrorResponse create instance of errorResponse
func NewErrorResponse(w http.ResponseWriter, statusCode int, description string) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(errorResponse{
		OK:          false,
		ErrorCode:   statusCode,
		Description: description,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
