package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

func WriteSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	WriteJSONResponse(w, http.StatusOK, data)
}

// SendJSONResponse sends a JSON response
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	WriteJSONResponse(w, statusCode, data)
}

// SendErrorResponse sends an error response
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	WriteErrorResponse(w, statusCode, message)
}


