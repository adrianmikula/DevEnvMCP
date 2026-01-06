package license

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	storage := NewStorage()
	assert.NotNil(t, storage)
	assert.NotEmpty(t, storage.configDir)
}

func TestSaveAndLoadLicense(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &Storage{
		configDir: tmpDir,
	}

	licenseKey := "test-license-key-12345"

	// Save license
	err := storage.SaveLicense(licenseKey)
	require.NoError(t, err)

	// Load license
	loaded, err := storage.LoadLicense()
	require.NoError(t, err)
	assert.Equal(t, licenseKey, loaded)
}

func TestLoadLicense_FileNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &Storage{
		configDir: tmpDir,
	}

	// Load from non-existent file
	loaded, err := storage.LoadLicense()
	require.NoError(t, err)
	assert.Empty(t, loaded) // Should return empty string, not error
}

func TestClearLicense(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &Storage{
		configDir: tmpDir,
	}

	// Save a license first
	err := storage.SaveLicense("test-key")
	require.NoError(t, err)

	// Verify it exists
	licenseFile := filepath.Join(tmpDir, "license.json")
	_, err = os.Stat(licenseFile)
	require.NoError(t, err)

	// Clear it
	err = storage.ClearLicense()
	require.NoError(t, err)

	// Verify it's gone
	_, err = os.Stat(licenseFile)
	assert.True(t, os.IsNotExist(err))
}

func TestClearLicense_FileNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &Storage{
		configDir: tmpDir,
	}

	// Clear non-existent license (should not error)
	err := storage.ClearLicense()
	require.NoError(t, err)
}

func TestSaveLicense_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "non-existent", "subdir")
	storage := &Storage{
		configDir: subDir,
	}

	// Save should create the directory
	err := storage.SaveLicense("test-key")
	require.NoError(t, err)

	// Verify directory was created
	_, err = os.Stat(subDir)
	assert.NoError(t, err)
}

