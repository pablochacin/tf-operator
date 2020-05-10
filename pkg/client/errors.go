package client

import (
	"errors"
	"fmt"
)

type ErrorReason string

const (
	// The object references in an operation doesn't exists
	ErrorReasonNotFound ErrorReason = "Not Found"
)

// TFOError extends error with a ErrorReason
type TFOError struct {
	error
	Reason ErrorReason
}

// NewFTOError creates a TFOError with an error description and a ErrorReason
func NewTFOError(desc string, reason ErrorReason) error {
	return TFOError{
		errors.New(desc),
		reason,
	}
}

// NewNotFoundError returns an error with NotFound cause
func NewNotFoundError(name string, class string, ns string) error {
	msg := fmt.Sprintf(" <%s> '%s' not found in <amespace> '%s'", class, name, ns)
	return NewTFOError(msg, ErrorReasonNotFound)
}

// IsNotFoundError indicates if the error is a not found error
func IsNotFoundError(err error) bool {
	var ftoErr TFOError
	if errors.As(err, &ftoErr) {
		return ftoErr.Reason == ErrorReasonNotFound
	}
	return false
}
