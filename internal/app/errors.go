package app

import "errors"

// Common errors used throughout the app package.
var (
	ErrNotImplemented = errors.New("not yet implemented")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrInvalidInput = errors.New("invalid input")
)
