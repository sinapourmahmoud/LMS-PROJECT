// Package errorhandler defines the standard error response shape used across the LMS backend.
package errorhandler

// Response represents an error response in the system.
type Response struct {
	Code       int         `json:"code"`
	Message    interface{} `json:"message"`
	StatusCode int         `json:"status_code"`
}

// New builds a Response from an error code, a message (a string, or e.g. a slice
// of validation field errors), and the HTTP status code it should be sent with.
func New(code int, message interface{}, statusCode int) *Response {
	return &Response{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}
