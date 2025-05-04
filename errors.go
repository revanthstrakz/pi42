package pi42

import "fmt"

// APIError represents an error returned by the Pi42 API
type APIError struct {
	StatusCode int    `json:"-"`
	ErrorCode  int    `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return fmt.Sprintf("API Error (Code: %d, Status: %d): %s", e.ErrorCode, e.StatusCode, e.Message)
}

// RequestError represents an error that occurs during API request
type RequestError struct {
	Message string
}

// Error implements the error interface
func (e RequestError) Error() string {
	return fmt.Sprintf("Request Error: %s", e.Message)
}
