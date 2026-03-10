// Package errors defines structured error types and exit codes for the MYN CLI.
// Exit codes per Spec §13.1.
package errors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Exit codes per Spec §13.1.
const (
	ExitSuccess    = 0 // Success
	ExitGeneral    = 1 // General error
	ExitUsage      = 2 // Usage error (bad flags, missing args)
	ExitAuth       = 3 // Authentication error
	ExitNetwork    = 4 // Network error (cannot reach backend)
	ExitAPI        = 5 // API error (4xx/5xx from backend)
	ExitRateLimit  = 6 // Rate limited (429)
)

// CLIError is a structured error with an exit code and optional hint.
type CLIError struct {
	Message  string
	Code     int
	Hint     string
	Cause    error
}

func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *CLIError) Unwrap() error {
	return e.Cause
}

// ExitCode returns the exit code for e.
// If e is nil or not a *CLIError, returns ExitGeneral.
func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	if errors.Is(err, context.Canceled) {
		return 130 // SIGINT standard
	}
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr.Code
	}
	return ExitGeneral
}

// New creates a CLIError with the given code and message.
func New(code int, message string) *CLIError {
	return &CLIError{Message: message, Code: code}
}

// Wrap wraps an existing error with a code and message.
func Wrap(code int, message string, cause error) *CLIError {
	return &CLIError{Message: message, Code: code, Cause: cause}
}

// WithHint attaches a hint message to the error.
func (e *CLIError) WithHint(hint string) *CLIError {
	e.Hint = hint
	return e
}

// Auth creates an authentication error (exit code 3).
func Auth(msg string, cause error) *CLIError {
	return &CLIError{
		Message: msg,
		Code:    ExitAuth,
		Hint:    "Run 'mynow login' to authenticate",
		Cause:   cause,
	}
}

// Network creates a network error (exit code 4).
func Network(cause error) *CLIError {
	return &CLIError{
		Message: "cannot reach the MYN backend",
		Code:    ExitNetwork,
		Hint:    "Check your internet connection and try again",
		Cause:   cause,
	}
}

// API creates an API error (exit code 5).
func API(statusCode int, body string) *CLIError {
	return &CLIError{
		Message: fmt.Sprintf("API error %d: %s", statusCode, strings.TrimSpace(body)),
		Code:    ExitAPI,
	}
}

// RateLimit creates a rate-limit error (exit code 6).
func RateLimit() *CLIError {
	return &CLIError{
		Message: "rate limited by the MYN backend",
		Code:    ExitRateLimit,
		Hint:    "Wait a moment and try again",
	}
}

// Usage creates a usage error (exit code 2).
func Usage(msg string) *CLIError {
	return &CLIError{Message: msg, Code: ExitUsage}
}

// JSONError serializes a CLIError for --json output per Spec §13.2.
func JSONError(err error) ([]byte, error) {
	var cliErr *CLIError
	if !errors.As(err, &cliErr) {
		cliErr = &CLIError{Message: err.Error(), Code: ExitGeneral}
	}
	type jsonErr struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
		Hint  string `json:"hint,omitempty"`
	}
	return json.Marshal(jsonErr{
		Error: cliErr.Message,
		Code:  cliErr.Code,
		Hint:  cliErr.Hint,
	})
}

// IsNetwork returns true if err is a network error.
func IsNetwork(err error) bool {
	var cliErr *CLIError
	return errors.As(err, &cliErr) && cliErr.Code == ExitNetwork
}

// IsAuth returns true if err is an auth error.
func IsAuth(err error) bool {
	var cliErr *CLIError
	return errors.As(err, &cliErr) && cliErr.Code == ExitAuth
}

// IsRateLimit returns true if err is a rate-limit error.
func IsRateLimit(err error) bool {
	var cliErr *CLIError
	return errors.As(err, &cliErr) && cliErr.Code == ExitRateLimit
}
