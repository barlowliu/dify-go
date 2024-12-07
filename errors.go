package dify-go

import "fmt"

// APIError represents an error returned by the Dify API.
type APIError struct {
    StatusCode int    `json:"status_code"`
    Code       string `json:"code"`
    Message    string `json:"message"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
    return fmt.Sprintf("APIError: %s - %s (status code: %d)", e.Code, e.Message, e.StatusCode)
}

