package http

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	OK          bool              `json:"ok"`
	ErrorCode   int               `json:"error_code,omitempty"`
	Description string            `json:"description,omitempty"`
}

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
