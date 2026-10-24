package database

import (
	"errors"
	"fmt"
)

// error categories; caller can map these to response.
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
	ErrInvalid  = errors.New("invalid input")
)

// specific errors, each wrapping a category so errors.Is matches both
var (
	ErrMemberNotFound     = fmt.Errorf("member %w", ErrNotFound)
	ErrInsufficientPoints = fmt.Errorf("insufficient points: %w", ErrInvalid)
	ErrNegativePoints     = fmt.Errorf("negative points passed: %w", ErrInvalid)
	ErrInvalidPointType   = fmt.Errorf("invalid point type: %w", ErrInvalid)
	ErrDuplicateEmail     = fmt.Errorf("email id already taken: %w", ErrConflict)
)
