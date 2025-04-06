package errs

import (
	"fmt"
)

type ErrorType string

const (
	ErrorTypeNotFound ErrorType = "not_found"

	ErrorTypeValidation ErrorType = "validation"

	ErrorTypeInternal ErrorType = "internal"

	ErrorTypeUnauthorized ErrorType = "unauthorized"
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}
