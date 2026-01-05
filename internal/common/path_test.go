package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		setup    func() func()
		validate func(t *testing.T, result string, err error)
	}{
		{
			name: "absolute path",
			path: "/tmp/test",
			validate: func(t *testing.T, result string, err error) {
				require.NoError(t, err)
				// On Windows, paths are converted, so just check it's absolute
				assert.True(t, filepath.IsAbs(result))
			},
		},
		{
			name: "relative path",
			path: "test",
			validate: func(t *testing.T, result string, err error) {
				require.NoError(t, err)
				abs, _ := filepath.Abs("test")
				assert.Equal(t, abs, result)
			},
		},
		{
			name: "with environment variable",
			path: "${HOME}/test",
			setup: func() func() {
				original := os.Getenv("HOME")
				os.Setenv("HOME", "/home/user")
				return func() {
					if original != "" {
						os.Setenv("HOME", original)
					}
				}
			},
			validate: func(t *testing.T, result string, err error) {
				require.NoError(t, err)
				// On Windows, paths use backslashes, so check for the expanded value
				assert.Contains(t, result, "home")
				assert.Contains(t, result, "user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cleanup func()
			if tt.setup != nil {
				cleanup = tt.setup()
				if cleanup != nil {
					defer cleanup()
				}
			}

			result, err := ResolvePath(tt.path)
			tt.validate(t, result, err)
		})
	}
}

func TestExpandPattern(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		setup    func() func()
		expected string
	}{
		{
			name:     "no variables",
			pattern:  "*.txt",
			expected: "*.txt",
		},
		{
			name:    "with environment variable",
			pattern: "${HOME}/test/*.txt",
			setup: func() func() {
				original := os.Getenv("HOME")
				os.Setenv("HOME", "/home/user")
				return func() {
					if original != "" {
						os.Setenv("HOME", original)
					}
				}
			},
			expected: "/home/user/test/*.txt",
		},
		{
			name:     "multiple variables",
			pattern:  "${VAR1}/${VAR2}/*.txt",
			setup: func() func() {
				os.Setenv("VAR1", "path1")
				os.Setenv("VAR2", "path2")
				return func() {
					os.Unsetenv("VAR1")
					os.Unsetenv("VAR2")
				}
			},
			expected: "path1/path2/*.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cleanup func()
			if tt.setup != nil {
				cleanup = tt.setup()
				if cleanup != nil {
					defer cleanup()
				}
			}

			result := ExpandPattern(tt.pattern)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinPaths(t *testing.T) {
	tests := []struct {
		name     string
		elements []string
		expected string
	}{
		{
			name:     "two elements",
			elements: []string{"path", "to", "file"},
			expected: filepath.Join("path", "to", "file"),
		},
		{
			name:     "single element",
			elements: []string{"path"},
			expected: "path",
		},
		{
			name:     "empty elements",
			elements: []string{},
			expected: "",
		},
		{
			name:     "with absolute path",
			elements: []string{"/absolute", "relative"},
			expected: filepath.Join("/absolute", "relative"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinPaths(tt.elements...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		validate func(t *testing.T, result string)
	}{
		{
			name: "clean path",
			path: "/path/to/file",
			validate: func(t *testing.T, result string) {
				expected := filepath.Clean("/path/to/file")
				assert.Equal(t, expected, result)
			},
		},
		{
			name: "with redundant separators",
			path: "/path//to///file",
			validate: func(t *testing.T, result string) {
				expected := filepath.Clean("/path//to///file")
				assert.Equal(t, expected, result)
			},
		},
		{
			name: "with .",
			path: "/path/./to/./file",
			validate: func(t *testing.T, result string) {
				expected := filepath.Clean("/path/./to/./file")
				assert.Equal(t, expected, result)
			},
		},
		{
			name: "with ..",
			path: "/path/to/../file",
			validate: func(t *testing.T, result string) {
				expected := filepath.Clean("/path/to/../file")
				assert.Equal(t, expected, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePath(tt.path)
			tt.validate(t, result)
		})
	}
}

func TestIsSubpath(t *testing.T) {
	tests := []struct {
		name     string
		basepath string
		subpath  string
		expected bool
		wantErr  bool
	}{
		{
			name:     "valid subpath",
			basepath: "/base",
			subpath:  "/base/sub",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "same path",
			basepath: "/base",
			subpath:  "/base",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "not a subpath",
			basepath: "/base",
			subpath:  "/other",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "directory traversal attempt",
			basepath: "/base",
			subpath:  "/base/../other",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "nested subpath",
			basepath: "/base",
			subpath:  "/base/sub/deep",
			expected: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use temp directory for actual path resolution
			tmpDir := t.TempDir()
			base := filepath.Join(tmpDir, "base")
			err := os.MkdirAll(base, 0755)
			require.NoError(t, err)

			var sub string
			if tt.name == "valid subpath" || tt.name == "nested subpath" {
				sub = filepath.Join(base, "sub")
				if tt.name == "nested subpath" {
					sub = filepath.Join(sub, "deep")
				}
				err = os.MkdirAll(sub, 0755)
				require.NoError(t, err)
			} else if tt.name == "same path" {
				sub = base
			} else if tt.name == "not a subpath" {
				sub = filepath.Join(tmpDir, "other")
				err = os.MkdirAll(sub, 0755)
				require.NoError(t, err)
			} else {
				// For directory traversal test, use actual relative path
				sub = filepath.Join(base, "..", "other")
			}

			result, err := IsSubpath(base, sub)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

