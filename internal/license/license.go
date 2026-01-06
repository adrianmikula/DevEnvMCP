package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

// License represents a license key and its validation state
type License struct {
	Key       string
	IsValid   bool
	Tier      string // "free", "pro", "enterprise"
	ExpiresAt *time.Time
	Features  []string
}

// LicenseValidator validates license keys
type LicenseValidator struct {
	secretKey string // Secret key for HMAC validation
}

// NewLicenseValidator creates a new license validator
// In production, the secret key should be embedded or fetched from a secure source
func NewLicenseValidator() *LicenseValidator {
	// Check for secret key in environment (for Apify/cloud deployments)
	secretKey := os.Getenv("SENTINEL_LICENSE_SECRET")
	if secretKey == "" {
		// Default secret for development (should be changed in production)
		secretKey = "dev-secret-key-change-in-production"
	}
	return &LicenseValidator{
		secretKey: secretKey,
	}
}

// ValidateLicense validates a license key
func (lv *LicenseValidator) ValidateLicense(key string) (*License, error) {
	if key == "" {
		return &License{
			Key:      "",
			IsValid:  false,
			Tier:     "free",
			Features: []string{"verify_build_freshness", "check_infrastructure_parity", "env_var_audit"},
		}, nil
	}

	// Check for Apify token (special format: apify_xxx)
	if strings.HasPrefix(key, "apify_") {
		return lv.validateApifyToken(key)
	}

	// Validate standard license key format: tier-hmac-timestamp
	parts := strings.Split(key, "-")
	if len(parts) != 3 {
		return &License{
			Key:     key,
			IsValid: false,
			Tier:    "free",
		}, fmt.Errorf("invalid license key format")
	}

	tier := parts[0]
	timestamp := parts[2]
	providedHMAC := parts[1]

	// Verify HMAC
	expectedHMAC := lv.computeHMAC(fmt.Sprintf("%s-%s", tier, timestamp))
	if !hmac.Equal([]byte(providedHMAC), []byte(expectedHMAC)) {
		return &License{
			Key:     key,
			IsValid: false,
			Tier:    "free",
		}, fmt.Errorf("invalid license key")
	}

	// Parse timestamp to check expiration
	var expiresAt *time.Time
	if timestamp != "lifetime" {
		expTime, err := time.Parse("20060102", timestamp)
		if err == nil {
			if time.Now().After(expTime) {
				return &License{
					Key:      key,
					IsValid:  false,
					Tier:     tier,
					ExpiresAt: &expTime,
				}, fmt.Errorf("license expired")
			}
			expiresAt = &expTime
		}
	}

	// Determine features based on tier
	features := lv.getFeaturesForTier(tier)

	return &License{
		Key:       key,
		IsValid:   true,
		Tier:      tier,
		ExpiresAt: expiresAt,
		Features:  features,
	}, nil
}

// validateApifyToken validates an Apify token
func (lv *LicenseValidator) validateApifyToken(token string) (*License, error) {
	// In Apify deployments, the token is validated by Apify's infrastructure
	// We just need to check if it's present and has the right format
	if strings.HasPrefix(token, "apify_") && len(token) > 10 {
		return &License{
			Key:      token,
			IsValid:  true,
			Tier:     "pro",
			Features: lv.getFeaturesForTier("pro"),
		}, nil
	}
	return &License{
		Key:     token,
		IsValid: false,
		Tier:    "free",
	}, fmt.Errorf("invalid Apify token")
}

// computeHMAC computes HMAC-SHA256 of the message
func (lv *LicenseValidator) computeHMAC(message string) string {
	h := hmac.New(sha256.New, []byte(lv.secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))[:16] // Use first 16 chars for shorter keys
}

// getFeaturesForTier returns the list of features for a given tier
func (lv *LicenseValidator) getFeaturesForTier(tier string) []string {
	switch tier {
	case "pro":
		return []string{
			"verify_build_freshness",
			"check_infrastructure_parity",
			"env_var_audit",
			"reconcile_environment", // Premium feature
			"auto_fix",
			"advanced_diagnostics",
		}
	case "enterprise":
		return []string{
			"verify_build_freshness",
			"check_infrastructure_parity",
			"env_var_audit",
			"reconcile_environment",
			"auto_fix",
			"advanced_diagnostics",
			"docker_orchestration",
			"priority_support",
			"custom_configs",
		}
	default: // free
		return []string{
			"verify_build_freshness",
			"check_infrastructure_parity",
			"env_var_audit",
		}
	}
}

// HasFeature checks if a license has a specific feature
func (l *License) HasFeature(feature string) bool {
	if !l.IsValid {
		return false
	}
	for _, f := range l.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// GetStripePaymentLink returns the Stripe payment link for Pro license
func GetStripePaymentLink() string {
	// This should be set via environment variable or config
	link := os.Getenv("SENTINEL_STRIPE_PAYMENT_LINK")
	if link == "" {
		// Default placeholder - should be replaced with actual Stripe link
		link = "https://buy.stripe.com/your-payment-link"
	}
	return link
}

// GetApifyActorURL returns the Apify Actor URL if deployed
func GetApifyActorURL() string {
	url := os.Getenv("SENTINEL_APIFY_ACTOR_URL")
	if url == "" {
		url = "https://api.apify.com/v2/actors/your-actor-id/run-sync"
	}
	return url
}

