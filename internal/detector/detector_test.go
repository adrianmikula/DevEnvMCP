package detector

import (
	"os"
	"path/filepath"
	"testing"

	"dev-env-sentinel/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectEcosystems(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (string, []*config.EcosystemConfig)
		expectedIDs []string
		wantErr     bool
	}{
		{
			name: "detects single ecosystem",
			setup: func(t *testing.T) (string, []*config.EcosystemConfig) {
				tmpDir := t.TempDir()
				pomPath := filepath.Join(tmpDir, "pom.xml")
				err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
				require.NoError(t, err)

				configs := []*config.EcosystemConfig{
					{
						Ecosystem: config.Ecosystem{
							ID: "java-maven",
							Detection: config.Detection{
								RequiredFiles: []string{"pom.xml"},
							},
						},
					},
				}
				return tmpDir, configs
			},
			expectedIDs: []string{"java-maven"},
			wantErr:     false,
		},
		{
			name: "detects multiple ecosystems",
			setup: func(t *testing.T) (string, []*config.EcosystemConfig) {
				tmpDir := t.TempDir()
				pomPath := filepath.Join(tmpDir, "pom.xml")
				err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
				require.NoError(t, err)

				packagePath := filepath.Join(tmpDir, "package.json")
				err = os.WriteFile(packagePath, []byte(`{"name": "test"}`), 0644)
				require.NoError(t, err)

				configs := []*config.EcosystemConfig{
					{
						Ecosystem: config.Ecosystem{
							ID: "java-maven",
							Detection: config.Detection{
								RequiredFiles: []string{"pom.xml"},
							},
						},
					},
					{
						Ecosystem: config.Ecosystem{
							ID: "npm",
							Detection: config.Detection{
								RequiredFiles: []string{"package.json"},
							},
						},
					},
				}
				return tmpDir, configs
			},
			expectedIDs: []string{"java-maven", "npm"},
			wantErr:     false,
		},
		{
			name: "no ecosystems detected",
			setup: func(t *testing.T) (string, []*config.EcosystemConfig) {
				tmpDir := t.TempDir()
				// Create a file that doesn't match any config
				otherPath := filepath.Join(tmpDir, "other.txt")
				err := os.WriteFile(otherPath, []byte("content"), 0644)
				require.NoError(t, err)

				configs := []*config.EcosystemConfig{
					{
						Ecosystem: config.Ecosystem{
							ID: "java-maven",
							Detection: config.Detection{
								RequiredFiles: []string{"pom.xml"},
							},
						},
					},
				}
				return tmpDir, configs
			},
			expectedIDs: []string{},
			wantErr:     false,
		},
		{
			name: "empty configs",
			setup: func(t *testing.T) (string, []*config.EcosystemConfig) {
				tmpDir := t.TempDir()
				return tmpDir, []*config.EcosystemConfig{}
			},
			expectedIDs: []string{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot, configs := tt.setup(t)

			ecosystems, err := DetectEcosystems(projectRoot, configs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, ecosystems, len(tt.expectedIDs))

			// Verify IDs
			ids := make(map[string]bool)
			for _, eco := range ecosystems {
				ids[eco.ID] = true
				assert.Equal(t, projectRoot, eco.ProjectRoot)
				assert.NotNil(t, eco.Config)
				assert.GreaterOrEqual(t, eco.Confidence, 0.5)
			}

			for _, expectedID := range tt.expectedIDs {
				assert.True(t, ids[expectedID], "expected ecosystem %s not found", expectedID)
			}
		})
	}
}

func TestIsEcosystemPresent(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) (string, *config.EcosystemConfig)
		expected   bool
		minConfidence float64
	}{
		{
			name: "all required files present",
			setup: func(t *testing.T) (string, *config.EcosystemConfig) {
				tmpDir := t.TempDir()
				pomPath := filepath.Join(tmpDir, "pom.xml")
				err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
				require.NoError(t, err)

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "java-maven",
						Detection: config.Detection{
							RequiredFiles: []string{"pom.xml"},
						},
					},
				}
				return tmpDir, cfg
			},
			expected: true,
			minConfidence: 0.5,
		},
		{
			name: "missing required file",
			setup: func(t *testing.T) (string, *config.EcosystemConfig) {
				tmpDir := t.TempDir()
				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "java-maven",
						Detection: config.Detection{
							RequiredFiles: []string{"pom.xml"},
						},
					},
				}
				return tmpDir, cfg
			},
			expected: false,
			minConfidence: 0.0,
		},
		{
			name: "partial required files",
			setup: func(t *testing.T) (string, *config.EcosystemConfig) {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "file1.txt")
				err := os.WriteFile(file1, []byte("content"), 0644)
				require.NoError(t, err)
				// file2.txt is missing

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "test",
						Detection: config.Detection{
							RequiredFiles: []string{"file1.txt", "file2.txt"},
						},
					},
				}
				return tmpDir, cfg
			},
			expected: false,
			minConfidence: 0.0,
		},
		{
			name: "with optional files boost confidence",
			setup: func(t *testing.T) (string, *config.EcosystemConfig) {
				tmpDir := t.TempDir()
				pomPath := filepath.Join(tmpDir, "pom.xml")
				err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
				require.NoError(t, err)

				mvnwPath := filepath.Join(tmpDir, "mvnw")
				err = os.WriteFile(mvnwPath, []byte("#!/bin/sh"), 0644)
				require.NoError(t, err)

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "java-maven",
						Detection: config.Detection{
							RequiredFiles: []string{"pom.xml"},
							OptionalFiles: []string{"mvnw"},
						},
					},
				}
				return tmpDir, cfg
			},
			expected: true,
			minConfidence: 0.5,
		},
		{
			name: "with directory patterns",
			setup: func(t *testing.T) (string, *config.EcosystemConfig) {
				tmpDir := t.TempDir()
				pomPath := filepath.Join(tmpDir, "pom.xml")
				err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
				require.NoError(t, err)

				srcDir := filepath.Join(tmpDir, "src", "main", "java")
				err = os.MkdirAll(srcDir, 0755)
				require.NoError(t, err)

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "java-maven",
						Detection: config.Detection{
							RequiredFiles:     []string{"pom.xml"},
							DirectoryPatterns: []string{"src/main/java"},
						},
					},
				}
				return tmpDir, cfg
			},
			expected: true,
			minConfidence: 0.5,
		},
		{
			name: "no required files but has optional",
			setup: func(t *testing.T) (string, *config.EcosystemConfig) {
				tmpDir := t.TempDir()
				mvnwPath := filepath.Join(tmpDir, "mvnw")
				err := os.WriteFile(mvnwPath, []byte("#!/bin/sh"), 0644)
				require.NoError(t, err)

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "java-maven",
						Detection: config.Detection{
							RequiredFiles: []string{},
							OptionalFiles: []string{"mvnw"},
						},
					},
				}
				return tmpDir, cfg
			},
			expected: true,
			minConfidence: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot, cfg := tt.setup(t)

			present, confidence := isEcosystemPresent(projectRoot, cfg)
			assert.Equal(t, tt.expected, present)
			if present {
				assert.GreaterOrEqual(t, confidence, tt.minConfidence)
			} else {
				assert.Less(t, confidence, 0.5)
			}
		})
	}
}

func TestDetectedEcosystem_Structure(t *testing.T) {
	tmpDir := t.TempDir()
	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
	require.NoError(t, err)

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "java-maven",
			Detection: config.Detection{
				RequiredFiles: []string{"pom.xml"},
			},
		},
	}

	ecosystems, err := DetectEcosystems(tmpDir, []*config.EcosystemConfig{cfg})
	require.NoError(t, err)
	require.Len(t, ecosystems, 1)

	eco := ecosystems[0]
	assert.Equal(t, "java-maven", eco.ID)
	assert.Equal(t, tmpDir, eco.ProjectRoot)
	assert.Equal(t, cfg, eco.Config)
	assert.GreaterOrEqual(t, eco.Confidence, 0.5)
	assert.LessOrEqual(t, eco.Confidence, 1.0)
}

