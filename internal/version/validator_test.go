package version

import (
	"testing"

	"dev-env-sentinel/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name     string
		info     *VersionInfo
		cfg      *config.EcosystemConfig
		expected bool
	}{
		{
			name: "valid version in range",
			info: &VersionInfo{
				Version: "17.0.0",
				Major:   "17",
				Minor:   "0",
				Patch:   "0",
			},
			cfg: &config.EcosystemConfig{
				Ecosystem: config.Ecosystem{
					Requirements: config.Requirements{
						MinVersion: "11",
						MaxVersion: "21",
					},
				},
			},
			expected: true,
		},
		{
			name: "version too old",
			info: &VersionInfo{
				Version: "8.0.0",
				Major:   "8",
			},
			cfg: &config.EcosystemConfig{
				Ecosystem: config.Ecosystem{
					Requirements: config.Requirements{
						MinVersion: "11",
					},
				},
			},
			expected: false,
		},
		{
			name: "version too new",
			info: &VersionInfo{
				Version: "22.0.0",
				Major:   "22",
			},
			cfg: &config.EcosystemConfig{
				Ecosystem: config.Ecosystem{
					Requirements: config.Requirements{
						MaxVersion: "21",
					},
				},
			},
			expected: false,
		},
		{
			name: "excluded version",
			info: &VersionInfo{
				Version: "15.0.0",
				Major:   "15",
			},
			cfg: &config.EcosystemConfig{
				Ecosystem: config.Ecosystem{
					Requirements: config.Requirements{
						ExcludedVersions: []string{"15"},
					},
				},
			},
			expected: false,
		},
		{
			name: "excluded runtime",
			info: &VersionInfo{
				Version: "17.0.0",
				RuntimeVariant: &RuntimeVariantInfo{
					Name:     "Oracle JDK",
					Provider: "Oracle",
				},
			},
			cfg: &config.EcosystemConfig{
				Ecosystem: config.Ecosystem{
					Requirements: config.Requirements{
						ExcludedRuntimes: []string{"Oracle JDK"},
					},
				},
			},
			expected: false,
		},
		{
			name: "preferred runtime",
			info: &VersionInfo{
				Version: "17.0.0",
				RuntimeVariant: &RuntimeVariantInfo{
					Name:     "Eclipse Temurin",
					Provider: "Adoptium",
				},
			},
			cfg: &config.EcosystemConfig{
				Ecosystem: config.Ecosystem{
					Requirements: config.Requirements{
						PreferredRuntimes: []string{"Eclipse Temurin"},
					},
				},
			},
			expected: true, // Preferred is valid, just not required
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateVersion(tt.info, tt.cfg)
			assert.Equal(t, tt.expected, result.IsValid)
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want int // -1: v1 < v2, 0: v1 == v2, 1: v1 > v2
	}{
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0", "1.0.0", 0},
		{"1.1.0", "1.0.0", 1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"17", "11", 1},
		{"17.0", "17.0.0", 0},
		{"1.2.3", "1.2.4", -1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestVersionGreaterOrEqual(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want bool
	}{
		{"2.0.0", "1.0.0", true},
		{"1.0.0", "2.0.0", false},
		{"1.0.0", "1.0.0", true},
		{"1.1.0", "1.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.v1+">="+tt.v2, func(t *testing.T) {
			assert.Equal(t, tt.want, versionGreaterOrEqual(tt.v1, tt.v2))
		})
	}
}

func TestVersionLessOrEqual(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want bool
	}{
		{"1.0.0", "2.0.0", true},
		{"2.0.0", "1.0.0", false},
		{"1.0.0", "1.0.0", true},
		{"1.0.0", "1.1.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"<="+tt.v2, func(t *testing.T) {
			assert.Equal(t, tt.want, versionLessOrEqual(tt.v1, tt.v2))
		})
	}
}

func TestGenerateSuggestions(t *testing.T) {
	info := &VersionInfo{
		Version:       "8.0.0",
		VersionManager: "sdkman",
	}

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			Requirements: config.Requirements{
				MinVersion:       "11",
				PreferredVersions: []string{"17", "21"},
			},
			VersionConfig: config.VersionConfig{
				VersionManagers: []config.VersionManager{
					{
						Name:           "sdkman",
						InstallCommand: "sdk install java {version}",
						SwitchCommand:  "sdk use java {version}",
					},
				},
			},
		},
	}

	issues := []ValidationIssue{
		{
			Type:     "version_too_old",
			Severity: "error",
			Message:  "Version 8.0.0 is below minimum required 11",
			Current:  "8.0.0",
			Required: "11",
		},
	}

	suggestions := generateSuggestions(info, cfg, issues)
	assert.NotEmpty(t, suggestions)
	assert.Equal(t, "switch_version", suggestions[0].Type)
	assert.Contains(t, suggestions[0].Versions, "17")
	assert.Contains(t, suggestions[0].Versions, "21")
}

