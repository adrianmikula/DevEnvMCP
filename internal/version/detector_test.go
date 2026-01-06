package version

import (
	"context"
	"runtime"
	"testing"
	"time"

	"dev-env-sentinel/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		pattern string
		want    *ParsedVersion
		wantErr bool
	}{
		{
			name:    "semantic version",
			output:  "java version \"17.0.9\"",
			pattern: `java version "([^"]+)"`,
			want: &ParsedVersion{
				Full:     "17.0.9",
				Semantic: "17.0.9",
				Major:    "17",
				Minor:    "0",
				Patch:    "9",
			},
			wantErr: false,
		},
		{
			name:    "version with suffix",
			output:  "Version 17.0.9+11",
			pattern: `Version ([^\s]+)`,
			want: &ParsedVersion{
				Full:     "17.0.9+11",
				Semantic: "17.0.9",
				Major:    "17",
				Minor:    "0",
				Patch:    "9",
			},
			wantErr: false,
		},
		{
			name:    "two-part version",
			output:  "Version 1.2",
			pattern: `Version (\d+\.\d+)`,
			want: &ParsedVersion{
				Full:     "1.2",
				Semantic: "1.2.0",
				Major:    "1",
				Minor:    "2",
				Patch:    "0",
			},
			wantErr: false,
		},
		{
			name:    "single version",
			output:  "Version 17",
			pattern: `Version (\d+)`,
			want: &ParsedVersion{
				Full:     "17",
				Semantic: "17.0.0",
				Major:    "17",
				Minor:    "0",
				Patch:    "0",
			},
			wantErr: false,
		},
		{
			name:    "no match",
			output:  "Some text",
			pattern: `Version (\d+)`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid pattern",
			output:  "Version 1.2.3",
			pattern: `[invalid`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseVersion(tt.output, tt.pattern)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.want.Full, result.Full)
				assert.Equal(t, tt.want.Semantic, result.Semantic)
				assert.Equal(t, tt.want.Major, result.Major)
				assert.Equal(t, tt.want.Minor, result.Minor)
				assert.Equal(t, tt.want.Patch, result.Patch)
			}
		})
	}
}

func TestDetectRuntimeVariant(t *testing.T) {
	versionCfg := config.VersionConfig{
		RuntimePattern: "(OpenJDK|Oracle|Eclipse Temurin)",
		RuntimeVariants: []config.RuntimeVariant{
			{
				Name:     "Eclipse Temurin",
				Provider: "Adoptium",
				Pattern:  "(?i)temurin",
			},
			{
				Name:     "OpenJDK",
				Provider: "OpenJDK",
				Pattern:  "(?i)openjdk",
			},
		},
	}

	tests := []struct {
		name    string
		output  string
		want    *RuntimeVariantInfo
		wantErr bool
	}{
		{
			name:   "OpenJDK detected",
			output: "openjdk version \"17.0.9\"",
			want: &RuntimeVariantInfo{
				Name:     "OpenJDK",
				Provider: "OpenJDK",
				FullName: "OpenJDK (OpenJDK)",
			},
			wantErr: false,
		},
		{
			name:   "Temurin detected",
			output: "openjdk version \"17.0.9\" 2023-10-17 LTS\nEclipse Temurin(TM) 64-Bit Server VM",
			want: &RuntimeVariantInfo{
				Name:     "Eclipse Temurin",
				Provider: "Adoptium",
				FullName: "Eclipse Temurin (Adoptium)",
			},
			wantErr: false,
		},
		{
			name:    "no match",
			output:  "Some other runtime",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detectRuntimeVariant(tt.output, versionCfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.want.Name, result.Name)
				assert.Equal(t, tt.want.Provider, result.Provider)
				assert.Equal(t, tt.want.FullName, result.FullName)
			}
		})
	}
}

func TestDetectVersionManager(t *testing.T) {
	versionCfg := config.VersionConfig{
		VersionManagers: []config.VersionManager{
			{
				Name:         "sdkman",
				CheckCommand: "command -v sdk",
			},
			{
				Name:         "jenv",
				CheckCommand: "command -v jenv",
			},
		},
	}

	ctx := context.Background()
	
	// This test depends on what's actually installed
	// Just verify the function doesn't crash
	manager := detectVersionManager(ctx, versionCfg)
	// Manager might be empty if none are installed, which is OK
	_ = manager
}

func TestDetectVersion(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows - requires sh")
	}

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			VersionConfig: config.VersionConfig{
				Language:       "java",
				VersionCommand: "echo 'openjdk version \"17.0.9\"'",
				VersionPattern: `openjdk version "([^"]+)"`,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := DetectVersion(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Equal(t, "java", info.Language)
	assert.Equal(t, "17.0.9", info.FullVersion)
	assert.Equal(t, "17.0.9", info.Version)
	assert.Equal(t, "17", info.Major)
	assert.Equal(t, "0", info.Minor)
	assert.Equal(t, "9", info.Patch)
}

func TestDetectVersion_CommandFails(t *testing.T) {
	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			VersionConfig: config.VersionConfig{
				Language:       "java",
				VersionCommand: "exit 1", // This will fail
				VersionPattern: `version "([^"]+)"`,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DetectVersion(ctx, cfg)
	assert.Error(t, err)
}

func TestDetectVersion_InvalidPattern(t *testing.T) {
	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			VersionConfig: config.VersionConfig{
				Language:       "java",
				VersionCommand: "echo 'version 1.2.3'",
				VersionPattern: `[invalid pattern`, // Invalid regex
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer cancel()

	_, err := DetectVersion(ctx, cfg)
	assert.Error(t, err)
}

