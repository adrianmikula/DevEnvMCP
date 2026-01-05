package version

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"dev-env-sentinel/internal/config"
)

// VersionInfo contains detected version information
type VersionInfo struct {
	Language      string
	Version       string
	FullVersion   string
	Major         string
	Minor         string
	Patch         string
	RuntimeVariant *RuntimeVariantInfo
	VersionManager string
	RawOutput     string
}

// RuntimeVariantInfo contains runtime variant information (Java-specific)
type RuntimeVariantInfo struct {
	Name     string
	Provider string
	FullName string
}

// DetectVersion detects the current language version
func DetectVersion(ctx context.Context, cfg *config.EcosystemConfig) (*VersionInfo, error) {
	versionCfg := cfg.Ecosystem.VersionConfig
	
	// Execute version command
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", versionCfg.VersionCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute version command: %w", err)
	}

	outputStr := strings.TrimSpace(string(output))
	
	// Parse version
	version, err := parseVersion(outputStr, versionCfg.VersionPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	info := &VersionInfo{
		Language:    versionCfg.Language,
		FullVersion: version.Full,
		Version:     version.Semantic,
		Major:       version.Major,
		Minor:       version.Minor,
		Patch:       version.Patch,
		RawOutput:   outputStr,
	}

	// Detect runtime variant if pattern provided
	if versionCfg.RuntimePattern != "" {
		runtime, err := detectRuntimeVariant(outputStr, versionCfg)
		if err == nil {
			info.RuntimeVariant = runtime
		}
	}

	// Detect version manager
	manager := detectVersionManager(ctx, versionCfg)
	if manager != "" {
		info.VersionManager = manager
	}

	return info, nil
}

// ParsedVersion contains parsed version components
type ParsedVersion struct {
	Full      string
	Semantic  string
	Major     string
	Minor     string
	Patch     string
}

// parseVersion parses version from output using pattern
func parseVersion(output, pattern string) (*ParsedVersion, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid version pattern: %w", err)
	}

	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return nil, fmt.Errorf("version pattern not found in output")
	}

	// Extract version components
	fullVersion := matches[1]
	
	// Parse semantic version (major.minor.patch)
	versionParts := strings.Split(fullVersion, ".")
	major := versionParts[0]
	minor := "0"
	patch := "0"
	
	if len(versionParts) > 1 {
		minor = versionParts[1]
	}
	if len(versionParts) > 2 {
		// Handle patch with possible suffix (e.g., "17.0.9+11")
		patch = strings.Split(versionParts[2], "+")[0]
		patch = strings.Split(patch, "-")[0]
	}

	semantic := fmt.Sprintf("%s.%s.%s", major, minor, patch)

	return &ParsedVersion{
		Full:     fullVersion,
		Semantic: semantic,
		Major:    major,
		Minor:    minor,
		Patch:    patch,
	}, nil
}

// detectRuntimeVariant detects runtime variant from output
func detectRuntimeVariant(output string, versionCfg config.VersionConfig) (*RuntimeVariantInfo, error) {
	if versionCfg.RuntimePattern == "" {
		return nil, fmt.Errorf("no runtime pattern configured")
	}

	// Try to match against known runtime variants
	for _, variant := range versionCfg.RuntimeVariants {
		re, err := regexp.Compile(variant.Pattern)
		if err != nil {
			continue
		}

		if re.MatchString(output) {
			return &RuntimeVariantInfo{
				Name:     variant.Name,
				Provider: variant.Provider,
				FullName: fmt.Sprintf("%s (%s)", variant.Name, variant.Provider),
			}, nil
		}
	}

	// Try generic runtime pattern if provided
	if versionCfg.RuntimePattern != "" {
		re, err := regexp.Compile(versionCfg.RuntimePattern)
		if err == nil {
			matches := re.FindStringSubmatch(output)
			if len(matches) > 1 {
				return &RuntimeVariantInfo{
					Name:     matches[1],
					Provider: "Unknown",
					FullName: matches[1],
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("runtime variant not detected")
}

// detectVersionManager detects which version manager is in use
func detectVersionManager(ctx context.Context, versionCfg config.VersionConfig) string {
	for _, manager := range versionCfg.VersionManagers {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		cmd := exec.CommandContext(ctx, "sh", "-c", manager.CheckCommand)
		err := cmd.Run()
		cancel()
		
		if err == nil {
			return manager.Name
		}
	}
	return ""
}

