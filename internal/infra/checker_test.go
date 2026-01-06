package infra

import (
	"context"
	"runtime"
	"testing"
	"time"

	"dev-env-sentinel/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckInfrastructure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows - requires sh")
	}

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "test",
			Infrastructure: config.Infrastructure{
				Services: []config.Service{
					{
						Name:         "test-service",
						Type:         "command",
						CheckCommand: "echo 'service running'",
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	report, err := CheckInfrastructure(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, report)

	assert.Len(t, report.Services, 1)
	assert.True(t, report.Services[0].Running)
	assert.True(t, report.Services[0].Healthy)
}

func TestCheckInfrastructure_ServiceFails(t *testing.T) {
	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "test",
			Infrastructure: config.Infrastructure{
				Services: []config.Service{
					{
						Name:         "failing-service",
						Type:         "command",
						CheckCommand: "exit 1", // This will fail
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	report, err := CheckInfrastructure(ctx, cfg)
	require.NoError(t, err)

	assert.Len(t, report.Services, 1)
	assert.False(t, report.Services[0].Running)
	assert.False(t, report.Services[0].Healthy)
}

func TestCheckInfrastructure_WithVersionExtract(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows - requires sh")
	}

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "test",
			Infrastructure: config.Infrastructure{
				Services: []config.Service{
					{
						Name:           "versioned-service",
						Type:           "command",
						CheckCommand:   "echo 'Version 1.2.3'",
						VersionExtract: "Version (\\d+\\.\\d+\\.\\d+)",
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	report, err := CheckInfrastructure(ctx, cfg)
	require.NoError(t, err)

	assert.Len(t, report.Services, 1)
	assert.Equal(t, "1.2.3", report.Services[0].Version)
}

func TestCheckService(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows - requires sh")
	}

	service := config.Service{
		Name:         "test-service",
		Type:         "command",
		CheckCommand: "echo 'running'",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := checkService(ctx, service)
	require.NoError(t, err)
	assert.True(t, status.Running)
	assert.True(t, status.Healthy)
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		pattern string
		want    string
		wantErr bool
	}{
		{
			name:    "valid version",
			output:  "Apache Maven 3.9.0",
			pattern: "Apache Maven (\\d+\\.\\d+\\.\\d+)",
			want:    "3.9.0",
			wantErr: false,
		},
		{
			name:    "no match",
			output:  "Some other text",
			pattern: "Version (\\d+\\.\\d+\\.\\d+)",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid pattern",
			output:  "Version 1.2.3",
			pattern: "[invalid regex",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := extractVersion(tt.output, tt.pattern)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, version)
			}
		})
	}
}

func TestCheckServiceHealth(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows - requires sh")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test successful command
	healthy, output, err := CheckServiceHealth(ctx, "echo 'healthy'", 2*time.Second)
	require.NoError(t, err)
	assert.True(t, healthy)
	assert.Contains(t, output, "healthy")

	// Test failing command
	healthy, _, err = CheckServiceHealth(ctx, "exit 1", 2*time.Second)
	assert.Error(t, err)
	assert.False(t, healthy)
}

func TestCheckServiceHealth_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Command that takes longer than timeout
	healthy, _, err := CheckServiceHealth(ctx, "sleep 5", 100*time.Millisecond)
	assert.Error(t, err)
	assert.False(t, healthy)
}

