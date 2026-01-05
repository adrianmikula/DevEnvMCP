package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNotFound_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrNotFound
		expected string
	}{
		{
			name: "file not found",
			err: &ErrNotFound{
				Resource: "file",
				Path:     "/path/to/file.txt",
			},
			expected: "file not found: /path/to/file.txt",
		},
		{
			name: "directory not found",
			err: &ErrNotFound{
				Resource: "directory",
				Path:     "/path/to/dir",
			},
			expected: "directory not found: /path/to/dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestErrInvalidConfig_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrInvalidConfig
		expected string
	}{
		{
			name: "with field",
			err: &ErrInvalidConfig{
				Field:   "ecosystem.id",
				Message: "required",
			},
			expected: "invalid config field ecosystem.id: required",
		},
		{
			name: "without field",
			err: &ErrInvalidConfig{
				Message: "failed to parse YAML",
			},
			expected: "invalid config: failed to parse YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestErrCommandFailed_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrCommandFailed
		expected string
	}{
		{
			name: "with output",
			err: &ErrCommandFailed{
				Command: "mvn clean",
				Output:  "BUILD FAILED",
			},
			expected: "command failed: mvn clean\noutput: BUILD FAILED",
		},
		{
			name: "with wrapped error",
			err: &ErrCommandFailed{
				Command: "mvn clean",
				Err:     errors.New("exit status 1"),
			},
			expected: "command failed: mvn clean: exit status 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestErrCommandFailed_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := &ErrCommandFailed{
		Command: "test",
		Err:     originalErr,
	}

	assert.Equal(t, originalErr, err.Unwrap())
	assert.True(t, errors.Is(err, originalErr))
}

