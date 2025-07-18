package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Details []ValidationError `json:"details"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	WriteJSON(w, http.StatusOK, response)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
	}
	WriteJSON(w, statusCode, response)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var details []ValidationError
		for _, fieldError := range validationErrors {
			details = append(details, ValidationError{
				Field:   fieldError.Field(),
				Message: getValidationMessage(fieldError),
			})
		}

		response := ValidationErrorResponse{
			Success: false,
			Error:   "Validation failed",
			Details: details,
		}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	WriteError(w, http.StatusBadRequest, err.Error())
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "min":
		return fieldError.Field() + " must be at least " + fieldError.Param() + " characters"
	case "max":
		return fieldError.Field() + " must be at most " + fieldError.Param() + " characters"
	case "email":
		return fieldError.Field() + " must be a valid email address"
	default:
		return fieldError.Field() + " is invalid"
	}
}
