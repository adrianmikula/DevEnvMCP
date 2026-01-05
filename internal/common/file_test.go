package common

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	info, err := GetFileInfo(testFile)
	require.NoError(t, err)
	assert.Equal(t, testFile, info.Path)
	assert.False(t, info.IsDir)
	assert.Greater(t, info.Size, int64(0))
	assert.False(t, info.ModTime.IsZero())
}

func TestGetFileInfo_NotFound(t *testing.T) {
	info, err := GetFileInfo("/nonexistent/file.txt")
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestGetFileInfo_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	info, err := GetFileInfo(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, tmpDir, info.Path)
	assert.True(t, info.IsDir)
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "exists.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	assert.True(t, FileExists(testFile))
	assert.False(t, FileExists(filepath.Join(tmpDir, "nonexistent.txt")))
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	assert.True(t, DirExists(tmpDir))
	assert.True(t, DirExists(subDir))
	assert.False(t, DirExists(filepath.Join(tmpDir, "nonexistent")))
}

func TestDirExists_FileNotDir(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "file.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	assert.False(t, DirExists(testFile))
}

func TestFindFilesByPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := []string{"file1.txt", "file2.txt", "file3.log"}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		err := os.WriteFile(path, []byte("content"), 0644)
		require.NoError(t, err)
	}

	// Create a subdirectory (should not be included)
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// Test pattern matching
	pattern := filepath.Join(tmpDir, "*.txt")
	matches, err := FindFilesByPattern(pattern)
	require.NoError(t, err)
	assert.Len(t, matches, 2)

	// Verify all matches are .txt files
	for _, match := range matches {
		assert.Contains(t, []string{"file1.txt", "file2.txt"}, filepath.Base(match))
	}
}

func TestFindFilesByPattern_NoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	pattern := filepath.Join(tmpDir, "*.nonexistent")
	matches, err := FindFilesByPattern(pattern)
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestFindDirsByPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test directories
	dirs := []string{"dir1", "dir2", "dir3"}
	for _, d := range dirs {
		path := filepath.Join(tmpDir, d)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err)
	}

	// Create a file (should not be included)
	testFile := filepath.Join(tmpDir, "file.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	// Test pattern matching
	pattern := filepath.Join(tmpDir, "dir*")
	matches, err := FindDirsByPattern(pattern)
	require.NoError(t, err)
	assert.Len(t, matches, 3)

	// Verify all matches are directories
	for _, match := range matches {
		assert.True(t, DirExists(match))
	}
}

func TestCompareTimestamps(t *testing.T) {
	tmpDir := t.TempDir()

	// Create first file
	file1 := filepath.Join(tmpDir, "file1.txt")
	err := os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	// Create second file (newer)
	file2 := filepath.Join(tmpDir, "file2.txt")
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)

	// Test comparison
	isNewer, err := CompareTimestamps(file1, file2)
	require.NoError(t, err)
	assert.False(t, isNewer) // file1 is older than file2

	isNewer, err = CompareTimestamps(file2, file1)
	require.NoError(t, err)
	assert.True(t, isNewer) // file2 is newer than file1
}

func TestCompareTimestamps_SameFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "file.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	isNewer, err := CompareTimestamps(testFile, testFile)
	require.NoError(t, err)
	assert.False(t, isNewer) // Same file, so not newer
}

func TestCompareTimestamps_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	err := os.WriteFile(existingFile, []byte("content"), 0644)
	require.NoError(t, err)

	nonexistent := filepath.Join(tmpDir, "nonexistent.txt")

	_, err = CompareTimestamps(existingFile, nonexistent)
	assert.Error(t, err)

	_, err = CompareTimestamps(nonexistent, existingFile)
	assert.Error(t, err)
}

