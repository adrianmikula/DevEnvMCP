package features

import (
	"fmt"

	"dev-env-sentinel/internal/license"
)

// FeatureManager manages feature flags and license-based feature access
type FeatureManager struct {
	license *license.License
}

// NewFeatureManager creates a new feature manager
func NewFeatureManager(lic *license.License) *FeatureManager {
	return &FeatureManager{
		license: lic,
	}
}

// IsEnabled checks if a feature is enabled for the current license
func (fm *FeatureManager) IsEnabled(feature string) bool {
	if fm.license == nil {
		return false
	}
	return fm.license.HasFeature(feature)
}

// RequireFeature returns an error if the feature is not available
func (fm *FeatureManager) RequireFeature(feature string) error {
	if !fm.IsEnabled(feature) {
		return &FeatureNotAvailableError{
			Feature: feature,
			Tier:    fm.license.Tier,
		}
	}
	return nil
}

// GetUpgradeMessage returns a message prompting the user to upgrade
func (fm *FeatureManager) GetUpgradeMessage(feature string) string {
	stripeLink := license.GetStripePaymentLink()
	apifyURL := license.GetApifyActorURL()
	
	return fmt.Sprintf(
		"The feature '%s' is only available in the Pro tier. "+
			"To unlock auto-fixes and advanced features, purchase a license:\n\n"+
			"• Stripe Payment Link: %s\n"+
			"• Apify Actor (Pay-Per-Event): %s\n\n"+
			"Once you have a license key, use the 'activate_pro' tool to activate it.",
		feature, stripeLink, apifyURL,
	)
}

// FeatureNotAvailableError is returned when a feature is not available
type FeatureNotAvailableError struct {
	Feature string
	Tier    string
}

func (e *FeatureNotAvailableError) Error() string {
	return fmt.Sprintf("feature '%s' is not available in tier '%s'", e.Feature, e.Tier)
}

