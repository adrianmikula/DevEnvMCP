package infra

import (
	"context"
	"fmt"

	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/version"
)

// CheckVersion checks language version and runtime compatibility
func CheckVersion(ctx context.Context, cfg *config.EcosystemConfig) (*VersionCheckResult, error) {
	// Detect current version
	versionInfo, err := version.DetectVersion(ctx, cfg)
	if err != nil {
		return &VersionCheckResult{
			Detected: false,
			Error:    err.Error(),
		}, nil
	}

	// Validate version
	validation := version.ValidateVersion(versionInfo, cfg)

	result := &VersionCheckResult{
		Detected:      true,
		VersionInfo:   versionInfo,
		IsValid:       validation.IsValid,
		Issues:        []string{},
		Suggestions:   []string{},
	}

	// Convert issues to strings
	for _, issue := range validation.Issues {
		result.Issues = append(result.Issues, issue.Message)
	}

	// Convert suggestions to strings
	for _, suggestion := range validation.Suggestions {
		msg := suggestion.Description
		if len(suggestion.Commands) > 0 {
			msg += fmt.Sprintf("\n  Commands:\n")
			for _, cmd := range suggestion.Commands {
				msg += fmt.Sprintf("    - %s\n", cmd)
			}
		}
		if len(suggestion.Versions) > 0 {
			msg += fmt.Sprintf("  Available versions: %v", suggestion.Versions)
		}
		result.Suggestions = append(result.Suggestions, msg)
	}

	return result, nil
}

// VersionCheckResult contains version check results
type VersionCheckResult struct {
	Detected    bool
	VersionInfo *version.VersionInfo
	IsValid     bool
	Issues      []string
	Suggestions []string
	Error       string
}

