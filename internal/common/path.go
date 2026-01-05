package common

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath expands environment variables and resolves to absolute path
func ResolvePath(path string) (string, error) {
	// Expand environment variables
	expanded := os.ExpandEnv(path)
	
	// Resolve to absolute path
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}
	
	return abs, nil
}

// ExpandPattern expands environment variables in a glob pattern
func ExpandPattern(pattern string) string {
	return os.ExpandEnv(pattern)
}

// JoinPaths joins path elements, handling both absolute and relative paths
func JoinPaths(elements ...string) string {
	return filepath.Join(elements...)
}

// NormalizePath normalizes a path (removes redundant separators, etc.)
func NormalizePath(path string) string {
	return filepath.Clean(path)
}

// IsSubpath checks if subpath is within basepath (prevents directory traversal)
func IsSubpath(basepath, subpath string) (bool, error) {
	baseAbs, err := filepath.Abs(basepath)
	if err != nil {
		return false, err
	}
	
	subAbs, err := filepath.Abs(subpath)
	if err != nil {
		return false, err
	}
	
	rel, err := filepath.Rel(baseAbs, subAbs)
	if err != nil {
		return false, err
	}
	
	return !strings.HasPrefix(rel, ".."), nil
}

