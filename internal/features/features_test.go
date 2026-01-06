package features

import (
	"testing"

	"dev-env-sentinel/internal/license"
	"github.com/stretchr/testify/assert"
)

func TestNewFeatureManager(t *testing.T) {
	lic := &license.License{
		IsValid: true,
		Tier:     "pro",
		Features: []string{"reconcile_environment"},
	}

	fm := NewFeatureManager(lic)
	assert.NotNil(t, fm)
	assert.Equal(t, lic, fm.license)
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		license  *license.License
		feature  string
		expected bool
	}{
		{
			name: "feature enabled",
			license: &license.License{
				IsValid: true,
				Tier:     "pro",
				Features: []string{"reconcile_environment"},
			},
			feature:  "reconcile_environment",
			expected: true,
		},
		{
			name: "feature not in list",
			license: &license.License{
				IsValid: true,
				Tier:     "pro",
				Features: []string{"other_feature"},
			},
			feature:  "reconcile_environment",
			expected: false,
		},
		{
			name:     "nil license",
			license:  nil,
			feature:  "any_feature",
			expected: false,
		},
		{
			name: "invalid license",
			license: &license.License{
				IsValid: false,
				Tier:     "free",
				Features: []string{"reconcile_environment"},
			},
			feature:  "reconcile_environment",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := NewFeatureManager(tt.license)
			assert.Equal(t, tt.expected, fm.IsEnabled(tt.feature))
		})
	}
}

func TestRequireFeature(t *testing.T) {
	tests := []struct {
		name    string
		license *license.License
		feature string
		wantErr bool
	}{
		{
			name: "feature available",
			license: &license.License{
				IsValid: true,
				Tier:     "pro",
				Features: []string{"reconcile_environment"},
			},
			feature: "reconcile_environment",
			wantErr: false,
		},
		{
			name: "feature not available",
			license: &license.License{
				IsValid: true,
				Tier:     "free",
				Features: []string{"other_feature"},
			},
			feature: "reconcile_environment",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := NewFeatureManager(tt.license)
			err := fm.RequireFeature(tt.feature)
			if tt.wantErr {
				assert.Error(t, err)
				var featureErr *FeatureNotAvailableError
				assert.ErrorAs(t, err, &featureErr)
				assert.Equal(t, tt.feature, featureErr.Feature)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUpgradeMessage(t *testing.T) {
	lic := &license.License{
		IsValid: true,
		Tier:     "free",
	}

	fm := NewFeatureManager(lic)
	msg := fm.GetUpgradeMessage("reconcile_environment")

	assert.Contains(t, msg, "reconcile_environment")
	assert.Contains(t, msg, "Pro tier")
	assert.Contains(t, msg, "Stripe Payment Link")
	assert.Contains(t, msg, "Apify Actor")
}

func TestFeatureNotAvailableError(t *testing.T) {
	err := &FeatureNotAvailableError{
		Feature: "reconcile_environment",
		Tier:    "free",
	}

	msg := err.Error()
	assert.Contains(t, msg, "reconcile_environment")
	assert.Contains(t, msg, "free")
}

