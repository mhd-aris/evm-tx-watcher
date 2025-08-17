package errors

import "fmt"

type ErrorCode string

const (
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeDatabase      ErrorCode = "DATABASE_ERROR"
	ErrCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest    ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeBlockchain    ErrorCode = "BLOCKCHAIN_ERROR"
	ErrCodeWebhook       ErrorCode = "WEBHOOK_ERROR"
)

type AppError struct {
	Code     ErrorCode              `json:"code"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	Internal error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code ErrorCode, message string, internal error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Internal: internal,
	}
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func NotFound(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func AlreadyExists(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeAlreadyExists,
		Message: fmt.Sprintf("%s already exists", resource),
	}
}

func ValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

func DatabaseError(internal error) *AppError {
	return &AppError{
		Code:     ErrCodeDatabase,
		Message:  "Database operation failed",
		Internal: internal,
	}
}

func InternalError(message string, internal error) *AppError {
	return &AppError{
		Code:     ErrCodeInternal,
		Message:  message,
		Internal: internal,
	}
}
