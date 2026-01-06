package license

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Storage handles license key persistence
type Storage struct {
	configDir string
}

// NewStorage creates a new license storage
func NewStorage() *Storage {
	// Use user's home directory for license storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	configDir := filepath.Join(homeDir, ".dev-env-sentinel")
	return &Storage{
		configDir: configDir,
	}
}

// SaveLicense saves a license key to disk
func (s *Storage) SaveLicense(key string) error {
	// Ensure config directory exists
	if err := os.MkdirAll(s.configDir, 0755); err != nil {
		return err
	}

	licenseFile := filepath.Join(s.configDir, "license.json")
	data := map[string]string{
		"key": key,
	}

	file, err := os.Create(licenseFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// LoadLicense loads a license key from disk
func (s *Storage) LoadLicense() (string, error) {
	licenseFile := filepath.Join(s.configDir, "license.json")
	
	data := make(map[string]string)
	file, err := os.Open(licenseFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No license file is OK
		}
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}

	return data["key"], nil
}

// ClearLicense removes the stored license
func (s *Storage) ClearLicense() error {
	licenseFile := filepath.Join(s.configDir, "license.json")
	return os.Remove(licenseFile)
}

