package http

import (
	"encoding/json"
	"net/http"
)

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"Pending":   true,
		"Completed": true,
		"Failed":    true,
		"Error":     true,
	}
	return validStatuses[status]
}

func sendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func sendSuccess(w http.ResponseWriter, code int, data interface{}, meta interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(SuccessResponse{Data: data, Meta: meta})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}
