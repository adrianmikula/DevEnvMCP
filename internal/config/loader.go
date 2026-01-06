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

// DiscoverEcosystemConfigs finds all ecosystem config files in the config directory structure
// New structure: config/languages/ (language yamls), config/languages/{lang}/ (tool yamls),
// config/infrastructure/ (infrastructure tools including databases, containers, docker, etc.)
// Falls back to old structure (language-configs, tool-configs) or baseDir for backwards compatibility
func DiscoverEcosystemConfigs(baseDir string) ([]*EcosystemConfig, error) {
	var configs []*EcosystemConfig

	configDir := filepath.Join(baseDir, "config")
	langDir := filepath.Join(configDir, "languages")
	infraDir := filepath.Join(configDir, "infrastructure")
	
	// Check if new structure exists
	if common.DirExists(configDir) {
		// Discover all configs in languages directory (language yamls and tool yamls in subdirs)
		if common.DirExists(langDir) {
			langConfigs, err := discoverConfigsInDir(langDir, true)
			if err != nil {
				return nil, fmt.Errorf("failed to discover language configs: %w", err)
			}
			configs = append(configs, langConfigs...)
		}

		// Discover infrastructure tool configs recursively (includes databases, containers, docker, etc.)
		if common.DirExists(infraDir) {
			infraConfigs, err := discoverConfigsInDir(infraDir, true)
			if err != nil {
				return nil, fmt.Errorf("failed to discover infrastructure configs: %w", err)
			}
			configs = append(configs, infraConfigs...)
		}
	} else {
		// Fallback to old structure: language-configs and tool-configs
		oldLangDir := filepath.Join(baseDir, "language-configs")
		oldToolDir := filepath.Join(baseDir, "tool-configs")
		
		if common.DirExists(oldLangDir) || common.DirExists(oldToolDir) {
			// Discover language configs
			if common.DirExists(oldLangDir) {
				langConfigs, err := discoverConfigsInDir(oldLangDir, false)
				if err != nil {
					return nil, fmt.Errorf("failed to discover language configs: %w", err)
				}
				configs = append(configs, langConfigs...)
			}

			// Discover tool configs recursively
			if common.DirExists(oldToolDir) {
				toolConfigs, err := discoverConfigsInDir(oldToolDir, true)
				if err != nil {
					return nil, fmt.Errorf("failed to discover tool configs: %w", err)
				}
				configs = append(configs, toolConfigs...)
			}
		} else {
			// Final fallback: discover configs directly in baseDir (for tests)
			if !common.DirExists(baseDir) {
				return nil, &common.ErrNotFound{Resource: "config directory", Path: baseDir}
			}
			
			flatConfigs, err := discoverConfigsInDir(baseDir, false)
			if err != nil {
				return nil, fmt.Errorf("failed to discover configs: %w", err)
			}
			configs = append(configs, flatConfigs...)
		}
	}

	return configs, nil
}

// discoverConfigsInDir finds all YAML config files in a directory, optionally recursing into subdirectories
// When recursive=false, it discovers YAML files in the current directory only
// When recursive=true, it discovers YAML files in the current directory AND recursively in subdirectories
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

		// Process YAML files in current directory
		if isYAMLFile(entry.Name()) {
			configPath := filepath.Join(dir, entry.Name())
			config, err := LoadEcosystemConfig(configPath)
			if err != nil {
				// Log error but continue with other configs
				continue
			}
			configs = append(configs, config)
		}
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

