package common

import (
	"os"
	"path/filepath"
	"time"
)

// FileInfo holds file information
type FileInfo struct {
	Path    string
	ModTime time.Time
	Size    int64
	IsDir   bool
}

// GetFileInfo returns file information
func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	return &FileInfo{
		Path:    path,
		ModTime: info.ModTime(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
	}, nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// FindFilesByPattern finds files matching a glob pattern
func FindFilesByPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil && !info.IsDir() {
			files = append(files, match)
		}
	}
	
	return files, nil
}

// FindDirsByPattern finds directories matching a glob pattern
func FindDirsByPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	
	var dirs []string
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil && info.IsDir() {
			dirs = append(dirs, match)
		}
	}
	
	return dirs, nil
}

// CompareTimestamps compares modification times of two files
// Returns true if first file is newer than second
func CompareTimestamps(file1, file2 string) (bool, error) {
	info1, err := GetFileInfo(file1)
	if err != nil {
		return false, err
	}
	
	info2, err := GetFileInfo(file2)
	if err != nil {
		return false, err
	}
	
	return info1.ModTime.After(info2.ModTime), nil
}

