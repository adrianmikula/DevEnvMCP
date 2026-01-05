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

// DiscoverEcosystemConfigs finds all ecosystem config files in a directory
func DiscoverEcosystemConfigs(configDir string) ([]*EcosystemConfig, error) {
	if !common.DirExists(configDir) {
		return nil, &common.ErrNotFound{Resource: "config directory", Path: configDir}
	}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var configs []*EcosystemConfig
	for _, entry := range entries {
		if entry.IsDir() || !isYAMLFile(entry.Name()) {
			continue
		}

		configPath := filepath.Join(configDir, entry.Name())
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

