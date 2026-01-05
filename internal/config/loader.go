package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dev-env-sentinel/internal/common"
	"gopkg.in/yaml.v3"
)

// LoadEcosystemConfig loads an ecosystem configuration from a YAML file
func LoadEcosystemConfig(path string) (*EcosystemConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &common.ErrNotFound{Resource: "config file", Path: path}
	}

	var config EcosystemConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, &common.ErrInvalidConfig{Message: fmt.Sprintf("failed to parse YAML: %v", err)}
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// DiscoverEcosystemConfigs finds all ecosystem config files in language-configs and tool-configs directories
// If those directories don't exist, it falls back to discovering configs in the base directory (for backwards compatibility)
func DiscoverEcosystemConfigs(baseDir string) ([]*EcosystemConfig, error) {
	var configs []*EcosystemConfig

	langDir := filepath.Join(baseDir, "language-configs")
	toolDir := filepath.Join(baseDir, "tool-configs")
	
	// Check if new structure exists
	if common.DirExists(langDir) || common.DirExists(toolDir) {
		// Discover language configs
		if common.DirExists(langDir) {
			langConfigs, err := discoverConfigsInDir(langDir, false)
			if err != nil {
				return nil, fmt.Errorf("failed to discover language configs: %w", err)
			}
			configs = append(configs, langConfigs...)
		}

		// Discover tool configs recursively
		if common.DirExists(toolDir) {
			toolConfigs, err := discoverConfigsInDir(toolDir, true)
			if err != nil {
				return nil, fmt.Errorf("failed to discover tool configs: %w", err)
			}
			configs = append(configs, toolConfigs...)
		}
	} else {
		// Fallback: discover configs directly in baseDir (for backwards compatibility and tests)
		if !common.DirExists(baseDir) {
			return nil, &common.ErrNotFound{Resource: "config directory", Path: baseDir}
		}
		
		flatConfigs, err := discoverConfigsInDir(baseDir, false)
		if err != nil {
			return nil, fmt.Errorf("failed to discover configs: %w", err)
		}
		configs = append(configs, flatConfigs...)
	}

	return configs, nil
}

// discoverConfigsInDir finds all YAML config files in a directory, optionally recursing into subdirectories
func discoverConfigsInDir(dir string, recursive bool) ([]*EcosystemConfig, error) {
	var configs []*EcosystemConfig

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if recursive {
				// Recursively discover configs in subdirectories
				subDir := filepath.Join(dir, entry.Name())
				subConfigs, err := discoverConfigsInDir(subDir, true)
				if err != nil {
					// Log error but continue with other directories
					continue
				}
				configs = append(configs, subConfigs...)
			}
			continue
		}

		if !isYAMLFile(entry.Name()) {
			continue
		}

		configPath := filepath.Join(dir, entry.Name())
		config, err := LoadEcosystemConfig(configPath)
		if err != nil {
			// Log error but continue with other configs
			continue
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// validateConfig validates the configuration structure
func validateConfig(config *EcosystemConfig) error {
	if config.Ecosystem.ID == "" {
		return &common.ErrInvalidConfig{Field: "ecosystem.id", Message: "required"}
	}

	if config.Ecosystem.Manifest.PrimaryFile == "" {
		return &common.ErrInvalidConfig{Field: "ecosystem.manifest.primary_file", Message: "required"}
	}

	return nil
}

// isYAMLFile checks if a file is a YAML file
func isYAMLFile(filename string) bool {
	ext := filepath.Ext(filename)
	extLower := strings.ToLower(ext)
	return extLower == ".yaml" || extLower == ".yml"
}

