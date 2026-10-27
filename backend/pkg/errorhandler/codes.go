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
	ForbiddenErrorCode             = 107
)

// User related error codes
const (
	UserNotFoundErrorCode = 200
)

// Enrollment related error codes
const (
	AlreadyEnrolledErrorCode = 400
	NotEnrolledErrorCode     = 404
)

const (
	CourseNotFoundErrorCode = 300
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
	ForbiddenMessage             = "You do not have permission to perform this action"
	AlreadyEnrolledErrorMessage  = "You are already enrolled in this course"
)

// User related error messages
const (
	UserNotFoundMessage = "User not found"
)

const (
	CourseNotFoundMessage = "Course not found"
)

const (
	NotEnrolledErrorMessage = "You are not enrolled in this course"
)
