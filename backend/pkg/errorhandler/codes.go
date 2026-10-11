// Package errorhandler defines the error codes and messages used across the LMS backend.
package errorhandler

// Generic error codes
const (
	InternalSystemErrorCode   = 1
	RequestBodyParseErrorCode = 2
	ValidationErrorCode       = 3
	DatabaseErrorCode         = 4
)

// Auth related error codes
const (
	InvalidCredentialsErrorCode    = 100
	EmailAlreadyExistsErrorCode    = 101
	UsernameAlreadyExistsErrorCode = 102
	InvalidAccessTokenErrorCode    = 103
	InvalidRefreshTokenErrorCode   = 104
	ExpiredAccessTokenErrorCode    = 105
	ExpiredRefreshTokenErrorCode   = 106
)

// User related error codes
const (
	UserNotFoundErrorCode = 200
)

// Generic error messages
const (
	InternalSystemErrorMessage   = "Internal system error"
	RequestBodyParseErrorMessage = "Invalid request body format"
	ValidationErrorMessage       = "Invalid request body"
	DatabaseErrorMessage         = "Database error"
)

// Auth related error messages
const (
	InvalidCredentialsMessage    = "Invalid credentials"
	EmailAlreadyExistsMessage    = "Email already exists"
	UsernameAlreadyExistsMessage = "Username already exists"
	InvalidAccessTokenMessage    = "Invalid access token"
	InvalidRefreshTokenMessage   = "Invalid refresh token"
	ExpiredAccessTokenMessage    = "Access token expired"
	ExpiredRefreshTokenMessage   = "Refresh token expired"
)

// User related error messages
const (
	UserNotFoundMessage = "User not found"
)
