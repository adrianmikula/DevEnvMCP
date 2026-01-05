package config

import (
	"os"
	"path/filepath"
	"testing"

	"dev-env-sentinel/internal/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEcosystemConfig(t *testing.T) {
	tests := []struct {
		name       string
		configYAML string
		wantErr    bool
		validate   func(t *testing.T, config *EcosystemConfig)
	}{
		{
			name: "valid config",
			configYAML: `
ecosystem:
  id: "test-ecosystem"
  name: "Test Ecosystem"
  manifest:
    primary_file: "pom.xml"
`,
			wantErr: false,
			validate: func(t *testing.T, config *EcosystemConfig) {
				assert.Equal(t, "test-ecosystem", config.Ecosystem.ID)
				assert.Equal(t, "Test Ecosystem", config.Ecosystem.Name)
				assert.Equal(t, "pom.xml", config.Ecosystem.Manifest.PrimaryFile)
			},
		},
		{
			name: "minimal valid config",
			configYAML: `
ecosystem:
  id: "minimal"
  manifest:
    primary_file: "package.json"
`,
			wantErr: false,
			validate: func(t *testing.T, config *EcosystemConfig) {
				assert.Equal(t, "minimal", config.Ecosystem.ID)
				assert.Equal(t, "package.json", config.Ecosystem.Manifest.PrimaryFile)
			},
		},
		{
			name:       "invalid YAML",
			configYAML: `invalid: yaml: [`,
			wantErr:    true,
		},
		{
			name:       "missing id",
			configYAML: `ecosystem: {}`,
			wantErr:    true,
		},
		{
			name: "missing primary_file",
			configYAML: `
ecosystem:
  id: "test"
`,
			wantErr: true,
		},
		{
			name: "full config structure",
			configYAML: `
ecosystem:
  id: "java-maven"
  name: "Java Maven"
  version: "1.0"
  detection:
    required_files:
      - "pom.xml"
    optional_files:
      - "mvnw"
  manifest:
    primary_file: "pom.xml"
    location: "."
    format: "xml"
  cache:
    locations:
      - "${HOME}/.m2/repository"
  build:
    output_directories:
      - "target/classes"
    clean_command: "mvn clean"
  verification:
    build_freshness:
      commands:
        - name: "check"
          type: "timestamp_compare"
          source: "pom.xml"
          target: "target/classes"
  environment:
    required_vars:
      - "JAVA_HOME"
  infrastructure:
    services:
      - name: "maven"
        type: "command"
        check_command: "mvn --version"
  reconciliation:
    fixes:
      - issue_type: "stale_build"
        command: "mvn clean"
        description: "Clean build"
`,
			wantErr: false,
			validate: func(t *testing.T, config *EcosystemConfig) {
				assert.Equal(t, "java-maven", config.Ecosystem.ID)
				assert.Len(t, config.Ecosystem.Detection.RequiredFiles, 1)
				assert.Len(t, config.Ecosystem.Cache.Locations, 1)
				assert.Len(t, config.Ecosystem.Verification.BuildFreshness.Commands, 1)
				assert.Len(t, config.Ecosystem.Reconciliation.Fixes, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "config.yaml")
			err := os.WriteFile(tmpFile, []byte(tt.configYAML), 0644)
			require.NoError(t, err)

			config, err := LoadEcosystemConfig(tmpFile)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
				// Verify error types
				if tt.name == "missing id" || tt.name == "missing primary_file" {
					var invalidConfigErr *common.ErrInvalidConfig
					assert.ErrorAs(t, err, &invalidConfigErr)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)
			if tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

func TestLoadEcosystemConfig_FileNotFound(t *testing.T) {
	config, err := LoadEcosystemConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
	assert.Nil(t, config)

	var notFoundErr *common.ErrNotFound
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestDiscoverEcosystemConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid config files
	validConfig1 := `
ecosystem:
  id: "ecosystem1"
  manifest:
    primary_file: "file1.xml"
`
	validConfig2 := `
ecosystem:
  id: "ecosystem2"
  manifest:
    primary_file: "file2.json"
`

	err := os.WriteFile(filepath.Join(tmpDir, "config1.yaml"), []byte(validConfig1), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "config2.yml"), []byte(validConfig2), 0644)
	require.NoError(t, err)

	// Create invalid config (should be skipped)
	invalidConfig := `invalid: yaml: [`
	err = os.WriteFile(filepath.Join(tmpDir, "invalid.yaml"), []byte(invalidConfig), 0644)
	require.NoError(t, err)

	// Create non-YAML file (should be skipped)
	err = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("not a config"), 0644)
	require.NoError(t, err)

	// Create subdirectory (should be skipped)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	configs, err := DiscoverEcosystemConfigs(tmpDir)
	require.NoError(t, err)
	assert.Len(t, configs, 2)

	// Verify IDs
	ids := make(map[string]bool)
	for _, cfg := range configs {
		ids[cfg.Ecosystem.ID] = true
	}
	assert.True(t, ids["ecosystem1"])
	assert.True(t, ids["ecosystem2"])
}

func TestDiscoverEcosystemConfigs_DirectoryNotFound(t *testing.T) {
	configs, err := DiscoverEcosystemConfigs("/nonexistent/directory")
	assert.Error(t, err)
	assert.Nil(t, configs)

	var notFoundErr *common.ErrNotFound
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestDiscoverEcosystemConfigs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configs, err := DiscoverEcosystemConfigs(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, configs)
}

func TestIsYAMLFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"yaml extension", "config.yaml", true},
		{"yml extension", "config.yml", true},
		{"txt extension", "config.txt", false},
		{"no extension", "config", false},
		{"json extension", "config.json", false},
		{"uppercase YAML", "CONFIG.YAML", true},
		{"uppercase YML", "CONFIG.YML", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isYAMLFile(tt.filename))
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *EcosystemConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &EcosystemConfig{
				Ecosystem: Ecosystem{
					ID: "test",
					Manifest: Manifest{
						PrimaryFile: "pom.xml",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing id",
			config: &EcosystemConfig{
				Ecosystem: Ecosystem{
					Manifest: Manifest{
						PrimaryFile: "pom.xml",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing primary_file",
			config: &EcosystemConfig{
				Ecosystem: Ecosystem{
					ID: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "empty id",
			config: &EcosystemConfig{
				Ecosystem: Ecosystem{
					ID: "",
					Manifest: Manifest{
						PrimaryFile: "pom.xml",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty primary_file",
			config: &EcosystemConfig{
				Ecosystem: Ecosystem{
					ID: "test",
					Manifest: Manifest{
						PrimaryFile: "",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				var invalidConfigErr *common.ErrInvalidConfig
				assert.ErrorAs(t, err, &invalidConfigErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
