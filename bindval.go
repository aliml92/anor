package anor

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrBinding indicates that binding *http.Request data to the receiver struct failed.
var ErrBinding = errors.New("binding error")

// ErrValidation indicates that the bound data did not pass validation.
var ErrValidation = errors.New("validation error")

// BindValidator is an interface that groups the basic Bind and Validate methods
type BindValidator interface {
	Binder
	Validator
}

// Binder is the interface that wraps the basic Bind method.
//
// Bind populates the receiver with data from the HTTP request.
type Binder interface {
	Bind(r *http.Request) error
}

// Validator is an interface that wraps the basic Validate method.
//
// Validate checks the receiver's data for validity.
type Validator interface {
	Validate() error
}

// BindValid binds HTTP request data to v and validates it.
// Returns ErrBinding or ErrValidation wrapped with details if either step fails.
func BindValid[T BindValidator](r *http.Request, v T) error {
	if err := v.Bind(r); err != nil {
		return fmt.Errorf("%w: %v", ErrBinding, err)
	}
	if err := v.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
