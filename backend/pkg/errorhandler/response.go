// Package errorhandler defines the standard error response shape used across the LMS backend.
package errorhandler

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

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



type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []FieldError {
	var out []FieldError

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, fe := range validationErrs {
			out = append(out, FieldError{
				Field:   fe.Field(),
				Message: fieldErrorMessage(fe),
			})
		}
	}

	return out
}

func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email address"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	default:
		return fe.Field() + " is invalid"
	}
}