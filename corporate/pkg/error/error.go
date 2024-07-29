package error

import (
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}

func NewErrorResponse(errors ...APIError) *ErrorResponse {
	return &ErrorResponse{
		Errors: errors,
	}
}

func BadRequest(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

func InternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

func NotFound(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}
