package verifier

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyBuildFreshness(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (string, *detector.DetectedEcosystem)
		expectIssue bool
		issueType   string
	}{
		{
			name: "detects stale build",
			setup: func(t *testing.T) (string, *detector.DetectedEcosystem) {
				tmpDir := t.TempDir()

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "test-ecosystem",
						Verification: config.Verification{
							BuildFreshness: config.BuildFreshness{
								Commands: []config.VerificationCommand{
									{
										Name:        "test_check",
										Type:        "timestamp_compare",
										Source:      "manifest.txt",
										Target:      "build/output.txt",
										Description: "Test timestamp comparison",
									},
								},
							},
						},
						Reconciliation: config.Reconciliation{
							Fixes: []config.Fix{
								{
									IssueType:   "stale_build",
									Command:     "build",
									Description: "Rebuild",
								},
							},
						},
					},
				}

				ecosystem := &detector.DetectedEcosystem{
					ID:          "test-ecosystem",
					Config:      cfg,
					Confidence:  1.0,
					ProjectRoot: tmpDir,
				}

				// Create older build output first
				buildDir := filepath.Join(tmpDir, "build")
				err := os.MkdirAll(buildDir, 0755)
				require.NoError(t, err)

				outputPath := filepath.Join(buildDir, "output.txt")
				err = os.WriteFile(outputPath, []byte("output"), 0644)
				require.NoError(t, err)

				// Wait a bit to ensure different timestamps
				time.Sleep(10 * time.Millisecond)

				// Create newer manifest file
				manifestPath := filepath.Join(tmpDir, "manifest.txt")
				err = os.WriteFile(manifestPath, []byte("manifest"), 0644)
				require.NoError(t, err)

				return tmpDir, ecosystem
			},
			expectIssue: true,
			issueType:   "stale_build",
		},
		{
			name: "healthy build - target newer than source",
			setup: func(t *testing.T) (string, *detector.DetectedEcosystem) {
				tmpDir := t.TempDir()

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "test-ecosystem",
						Verification: config.Verification{
							BuildFreshness: config.BuildFreshness{
								Commands: []config.VerificationCommand{
									{
										Name:        "test_check",
										Type:        "timestamp_compare",
										Source:      "manifest.txt",
										Target:      "build/output.txt",
										Description: "Test timestamp comparison",
									},
								},
							},
						},
					},
				}

				ecosystem := &detector.DetectedEcosystem{
					ID:          "test-ecosystem",
					Config:      cfg,
					Confidence:  1.0,
					ProjectRoot: tmpDir,
				}

				// Create manifest file
				manifestPath := filepath.Join(tmpDir, "manifest.txt")
				err := os.WriteFile(manifestPath, []byte("manifest"), 0644)
				require.NoError(t, err)

				// Wait a bit
				time.Sleep(10 * time.Millisecond)

				// Create newer build output
				buildDir := filepath.Join(tmpDir, "build")
				err = os.MkdirAll(buildDir, 0755)
				require.NoError(t, err)

				outputPath := filepath.Join(buildDir, "output.txt")
				err = os.WriteFile(outputPath, []byte("output"), 0644)
				require.NoError(t, err)

				return tmpDir, ecosystem
			},
			expectIssue: false,
		},
		{
			name: "missing target file",
			setup: func(t *testing.T) (string, *detector.DetectedEcosystem) {
				tmpDir := t.TempDir()

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "test-ecosystem",
						Verification: config.Verification{
							BuildFreshness: config.BuildFreshness{
								Commands: []config.VerificationCommand{
									{
										Name:        "test_check",
										Type:        "timestamp_compare",
										Source:      "manifest.txt",
										Target:      "build/output.txt",
										Description: "Test timestamp comparison",
									},
								},
							},
						},
					},
				}

				ecosystem := &detector.DetectedEcosystem{
					ID:          "test-ecosystem",
					Config:      cfg,
					Confidence:  1.0,
					ProjectRoot: tmpDir,
				}

				// Create manifest file only
				manifestPath := filepath.Join(tmpDir, "manifest.txt")
				err := os.WriteFile(manifestPath, []byte("manifest"), 0644)
				require.NoError(t, err)

				return tmpDir, ecosystem
			},
			expectIssue: true,
			issueType:   "missing_target",
		},
		{
			name: "target pattern matching",
			setup: func(t *testing.T) (string, *detector.DetectedEcosystem) {
				tmpDir := t.TempDir()

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "test-ecosystem",
						Verification: config.Verification{
							BuildFreshness: config.BuildFreshness{
								Commands: []config.VerificationCommand{
									{
										Name:          "test_check",
										Type:          "timestamp_compare",
										Source:        "manifest.txt",
										TargetPattern: filepath.Join("build", "*", "*.class"),
										Description:   "Test pattern matching",
									},
								},
							},
						},
						Reconciliation: config.Reconciliation{
							Fixes: []config.Fix{
								{
									IssueType:   "stale_build",
									Command:     "build",
									Description: "Rebuild",
								},
							},
						},
					},
				}

				ecosystem := &detector.DetectedEcosystem{
					ID:          "test-ecosystem",
					Config:      cfg,
					Confidence:  1.0,
					ProjectRoot: tmpDir,
				}

				// Create older build output first
				buildDir := filepath.Join(tmpDir, "build", "classes")
				err := os.MkdirAll(buildDir, 0755)
				require.NoError(t, err)

				classFile := filepath.Join(buildDir, "Test.class")
				err = os.WriteFile(classFile, []byte("class"), 0644)
				require.NoError(t, err)

				// Wait a bit to ensure different timestamps
				time.Sleep(10 * time.Millisecond)

				// Create newer manifest file
				manifestPath := filepath.Join(tmpDir, "manifest.txt")
				err = os.WriteFile(manifestPath, []byte("manifest"), 0644)
				require.NoError(t, err)

				return tmpDir, ecosystem
			},
			expectIssue: true,
			issueType:   "stale_build",
		},
		{
			name: "no verification commands",
			setup: func(t *testing.T) (string, *detector.DetectedEcosystem) {
				tmpDir := t.TempDir()

				cfg := &config.EcosystemConfig{
					Ecosystem: config.Ecosystem{
						ID: "test-ecosystem",
						Verification: config.Verification{
							BuildFreshness: config.BuildFreshness{
								Commands: []config.VerificationCommand{},
							},
						},
					},
				}

				ecosystem := &detector.DetectedEcosystem{
					ID:          "test-ecosystem",
					Config:      cfg,
					Confidence:  1.0,
					ProjectRoot: tmpDir,
				}

				return tmpDir, ecosystem
			},
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot, ecosystem := tt.setup(t)

			report, err := VerifyBuildFreshness(projectRoot, ecosystem)
			require.NoError(t, err)
			require.NotNil(t, report)
			assert.Equal(t, ecosystem.ID, report.EcosystemID)

			if tt.expectIssue {
				assert.False(t, report.IsHealthy)
				assert.NotEmpty(t, report.Issues)
				if tt.issueType != "" {
					found := false
					for _, issue := range report.Issues {
						if issue.Type == tt.issueType {
							found = true
							break
						}
					}
					assert.True(t, found, "expected issue type %s not found", tt.issueType)
				}
			} else {
				assert.True(t, report.IsHealthy)
				assert.Empty(t, report.Issues)
			}
		})
	}
}

func TestVerifyBuildFreshness_IssueDetails(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "test-ecosystem",
			Verification: config.Verification{
				BuildFreshness: config.BuildFreshness{
					Commands: []config.VerificationCommand{
						{
							Name:        "test_check",
							Type:        "timestamp_compare",
							Source:      "manifest.txt",
							Target:      "build/output.txt",
							Description: "Test timestamp comparison",
						},
					},
				},
			},
			Reconciliation: config.Reconciliation{
				Fixes: []config.Fix{
					{
						IssueType:   "stale_build",
						Command:     "mvn clean",
						Description: "Clean build",
					},
				},
			},
		},
	}

	ecosystem := &detector.DetectedEcosystem{
		ID:          "test-ecosystem",
		Config:      cfg,
		Confidence:  1.0,
		ProjectRoot: tmpDir,
	}

	// Create older build output first
	buildDir := filepath.Join(tmpDir, "build")
	err := os.MkdirAll(buildDir, 0755)
	require.NoError(t, err)

	outputPath := filepath.Join(buildDir, "output.txt")
	err = os.WriteFile(outputPath, []byte("output"), 0644)
	require.NoError(t, err)

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	// Create newer manifest file
	manifestPath := filepath.Join(tmpDir, "manifest.txt")
	err = os.WriteFile(manifestPath, []byte("manifest"), 0644)
	require.NoError(t, err)

	report, err := VerifyBuildFreshness(tmpDir, ecosystem)
	require.NoError(t, err)
	require.Len(t, report.Issues, 1)

	issue := report.Issues[0]
	assert.Equal(t, "stale_build", issue.Type)
	assert.Equal(t, "error", issue.Severity)
	assert.NotEmpty(t, issue.Message)
	assert.True(t, issue.FixAvailable)
	assert.Equal(t, "mvn clean", issue.FixCommand)
}

