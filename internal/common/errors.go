package common

import "fmt"

// Error types for better error handling and context

// ErrNotFound indicates a resource was not found
type ErrNotFound struct {
	Resource string
	Path     string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Resource, e.Path)
}

// ErrInvalidConfig indicates invalid configuration
type ErrInvalidConfig struct {
	Field   string
	Message string
}

func (e *ErrInvalidConfig) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("invalid config field %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("invalid config: %s", e.Message)
}

// ErrCommandFailed indicates a command execution failed
type ErrCommandFailed struct {
	Command string
	Output  string
	Err     error
}

func (e *ErrCommandFailed) Error() string {
	if e.Output != "" {
		return fmt.Sprintf("command failed: %s\noutput: %s", e.Command, e.Output)
	}
	return fmt.Sprintf("command failed: %s: %v", e.Command, e.Err)
}

func (e *ErrCommandFailed) Unwrap() error {
	return e.Err
}

