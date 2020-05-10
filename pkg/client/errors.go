package client

import (
	"errors"
	"fmt"
)

type ErrorReason string

const (
	// The object referenced in an operation doesn't exists
	ErrorReasonNotFound ErrorReason = "Not Found"

    // The object referenced in an operation already exists 
    ErrorReasonAlreadyExists ErrorReason = "Already Exists"
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

// NewAlreadyExistsError returns an error with AlreadyExists cause
func NewAlreadyExistsError(name string, class string, ns string) error {
	msg := fmt.Sprintf(" <%s> '%s' already exists in <amespace> '%s'", class, name, ns)
	return NewTFOError(msg, ErrorReasonAlreadyExists)
}

// IsNotFoundError indicates if the error has the given cause
func Is(err error, reason ErrorReason) bool {
	var ftoErr TFOError
	if errors.As(err, &ftoErr) {
		return ftoErr.Reason == reason
	}
	return false
}
