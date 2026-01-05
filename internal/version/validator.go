package version

import (
	"fmt"
	"strings"

	"dev-env-sentinel/internal/config"
)

// ValidationResult contains version validation results
type ValidationResult struct {
	IsValid        bool
	Issues         []ValidationIssue
	Suggestions    []Suggestion
}

// ValidationIssue describes a version compatibility issue
type ValidationIssue struct {
	Type        string // "version_too_old", "version_too_new", "wrong_runtime", etc.
	Severity    string // "error", "warning"
	Message     string
	Current     string
	Required    string
}

// Suggestion provides actionable fix suggestions
type Suggestion struct {
	Type        string // "install_version", "switch_version", "switch_runtime"
	Description string
	Commands    []string
	Versions    []string // Available versions that would work
}

// ValidateVersion validates version against requirements
func ValidateVersion(info *VersionInfo, cfg *config.EcosystemConfig) *ValidationResult {
	result := &ValidationResult{
		IsValid:     true,
		Issues:      []ValidationIssue{},
		Suggestions: []Suggestion{},
	}

	req := cfg.Ecosystem.Requirements

	// Check minimum version
	if req.MinVersion != "" {
		if !versionGreaterOrEqual(info.Version, req.MinVersion) {
			result.IsValid = false
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "version_too_old",
				Severity: "error",
				Message:  fmt.Sprintf("Version %s is below minimum required %s", info.Version, req.MinVersion),
				Current:  info.Version,
				Required: req.MinVersion,
			})
		}
	}

	// Check maximum version
	if req.MaxVersion != "" {
		if !versionLessOrEqual(info.Version, req.MaxVersion) {
			result.IsValid = false
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "version_too_new",
				Severity: "error",
				Message:  fmt.Sprintf("Version %s exceeds maximum allowed %s", info.Version, req.MaxVersion),
				Current:  info.Version,
				Required: req.MaxVersion,
			})
		}
	}

	// Check excluded versions
	for _, excluded := range req.ExcludedVersions {
		if info.Version == excluded || strings.HasPrefix(info.Version, excluded+".") {
			result.IsValid = false
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "version_excluded",
				Severity: "error",
				Message:  fmt.Sprintf("Version %s is excluded", info.Version),
				Current:  info.Version,
				Required: "different version",
			})
		}
	}

	// Check runtime variant (Java-specific)
	if info.RuntimeVariant != nil {
		// Check excluded runtimes
		for _, excluded := range req.ExcludedRuntimes {
			if info.RuntimeVariant.Name == excluded || info.RuntimeVariant.Provider == excluded {
				result.IsValid = false
				result.Issues = append(result.Issues, ValidationIssue{
					Type:     "runtime_excluded",
					Severity: "warning",
					Message:  fmt.Sprintf("Runtime %s is not recommended", info.RuntimeVariant.FullName),
					Current:  info.RuntimeVariant.FullName,
					Required: "different runtime",
				})
			}
		}

		// Check preferred runtimes
		if len(req.PreferredRuntimes) > 0 {
			isPreferred := false
			for _, preferred := range req.PreferredRuntimes {
				if info.RuntimeVariant.Name == preferred || info.RuntimeVariant.Provider == preferred {
					isPreferred = true
					break
				}
			}
			if !isPreferred {
				result.Issues = append(result.Issues, ValidationIssue{
					Type:     "runtime_not_preferred",
					Severity: "warning",
					Message:  fmt.Sprintf("Runtime %s is not in preferred list", info.RuntimeVariant.FullName),
					Current:  info.RuntimeVariant.FullName,
					Required: strings.Join(req.PreferredRuntimes, " or "),
				})
			}
		}
	}

	// Generate suggestions if there are issues
	if !result.IsValid {
		result.Suggestions = generateSuggestions(info, cfg, result.Issues)
	}

	return result
}

// versionGreaterOrEqual compares semantic versions
func versionGreaterOrEqual(v1, v2 string) bool {
	return compareVersions(v1, v2) >= 0
}

// versionLessOrEqual compares semantic versions
func versionLessOrEqual(v1, v2 string) bool {
	return compareVersions(v1, v2) <= 0
}

// compareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &p1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &p2)
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	return 0
}

// generateSuggestions generates fix suggestions
func generateSuggestions(info *VersionInfo, cfg *config.EcosystemConfig, issues []ValidationIssue) []Suggestion {
	var suggestions []Suggestion

	// Find version manager
	var manager *config.VersionManager
	for _, vm := range cfg.Ecosystem.VersionConfig.VersionManagers {
		if vm.Name == info.VersionManager {
			manager = &vm
			break
		}
	}

	// Generate suggestions based on issue types
	for _, issue := range issues {
		switch issue.Type {
		case "version_too_old", "version_too_new", "version_excluded":
			suggestion := Suggestion{
				Type:        "switch_version",
				Description: fmt.Sprintf("Switch to a compatible version (required: %s)", issue.Required),
			}

			// Add preferred versions if available
			if len(cfg.Ecosystem.Requirements.PreferredVersions) > 0 {
				suggestion.Versions = cfg.Ecosystem.Requirements.PreferredVersions
			} else {
				// Suggest versions in range
				if cfg.Ecosystem.Requirements.MinVersion != "" && cfg.Ecosystem.Requirements.MaxVersion != "" {
					suggestion.Versions = []string{
						cfg.Ecosystem.Requirements.MinVersion,
						cfg.Ecosystem.Requirements.MaxVersion,
					}
				}
			}

			// Add commands if version manager available
			if manager != nil {
				for _, version := range suggestion.Versions {
					installCmd := strings.ReplaceAll(manager.InstallCommand, "{version}", version)
					switchCmd := strings.ReplaceAll(manager.SwitchCommand, "{version}", version)
					suggestion.Commands = append(suggestion.Commands, installCmd, switchCmd)
				}
			}

			suggestions = append(suggestions, suggestion)

		case "runtime_excluded", "runtime_not_preferred":
			suggestion := Suggestion{
				Type:        "switch_runtime",
				Description: fmt.Sprintf("Switch to a preferred runtime: %s", issue.Required),
			}

			// Add preferred runtime options
			if len(cfg.Ecosystem.Requirements.PreferredRuntimes) > 0 {
				suggestion.Versions = cfg.Ecosystem.Requirements.PreferredRuntimes
			}

			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

