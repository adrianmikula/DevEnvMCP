package verifier

import (
	"fmt"
	"path/filepath"
	"time"

	"dev-env-sentinel/internal/common"
	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
)

// FreshnessReport contains the results of build freshness verification
type FreshnessReport struct {
	EcosystemID string
	IsHealthy   bool
	Issues      []Issue
}

// Issue represents a detected problem
type Issue struct {
	Type        string
	Severity    string
	Message     string
	FixAvailable bool
	FixCommand  string
}

// VerifyBuildFreshness verifies build freshness for a detected ecosystem
func VerifyBuildFreshness(projectRoot string, ecosystem *detector.DetectedEcosystem) (*FreshnessReport, error) {
	report := &FreshnessReport{
		EcosystemID: ecosystem.ID,
		IsHealthy:   true,
		Issues:      []Issue{},
	}

	cfg := ecosystem.Config
	verification := cfg.Ecosystem.Verification.BuildFreshness

	// Execute verification commands
	for _, cmd := range verification.Commands {
		issue, err := executeVerificationCommand(cmd, projectRoot, ecosystem)
		if err != nil {
			// Log error but continue with other checks
			continue
		}

		if issue != nil {
			report.IsHealthy = false
			report.Issues = append(report.Issues, *issue)
		}
	}

	return report, nil
}

// executeVerificationCommand executes a single verification command
func executeVerificationCommand(cmd config.VerificationCommand, projectRoot string, ecosystem *detector.DetectedEcosystem) (*Issue, error) {
	switch cmd.Type {
	case "timestamp_compare":
		return verifyTimestampCompare(cmd, projectRoot, ecosystem)
	case "command":
		return verifyCommand(cmd, projectRoot)
	default:
		return nil, fmt.Errorf("unknown verification command type: %s", cmd.Type)
	}
}

// verifyTimestampCompare verifies timestamp comparison
func verifyTimestampCompare(cmd config.VerificationCommand, projectRoot string, ecosystem *detector.DetectedEcosystem) (*Issue, error) {
	// Resolve source path
	sourcePath := filepath.Join(projectRoot, common.ExpandPattern(cmd.Source))
	if !common.FileExists(sourcePath) {
		return nil, fmt.Errorf("source file not found: %s", sourcePath)
	}

	sourceInfo, err := common.GetFileInfo(sourcePath)
	if err != nil {
		return nil, err
	}

	// Handle target pattern
	if cmd.TargetPattern != "" {
		return verifyTimestampPattern(sourceInfo, cmd.TargetPattern, projectRoot, cmd, ecosystem)
	}

	// Handle single target file
	if cmd.Target != "" {
		targetPath := filepath.Join(projectRoot, common.ExpandPattern(cmd.Target))
		if !common.FileExists(targetPath) {
			return &Issue{
				Type:        "missing_target",
				Severity:    "warning",
				Message:     fmt.Sprintf("Target file not found: %s", cmd.Target),
				FixAvailable: false,
			}, nil
		}

		targetInfo, err := common.GetFileInfo(targetPath)
		if err != nil {
			return nil, err
		}

		if sourceInfo.ModTime.After(targetInfo.ModTime) {
			return &Issue{
				Type:        "stale_build",
				Severity:    "error",
				Message:     fmt.Sprintf("%s is newer than %s", cmd.Source, cmd.Target),
				FixAvailable: true,
				FixCommand:  getFixCommand(ecosystem, "stale_build"),
			}, nil
		}
	}

	return nil, nil
}

// verifyTimestampPattern verifies timestamp against a pattern
func verifyTimestampPattern(sourceInfo *common.FileInfo, pattern string, projectRoot string, cmd config.VerificationCommand, ecosystem *detector.DetectedEcosystem) (*Issue, error) {
	expandedPattern := common.ExpandPattern(pattern)
	fullPattern := filepath.Join(projectRoot, expandedPattern)

	matches, err := common.FindFilesByPattern(fullPattern)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return &Issue{
			Type:        "missing_build_output",
			Severity:    "warning",
			Message:     fmt.Sprintf("No files found matching pattern: %s", pattern),
			FixAvailable: false,
		}, nil
	}

	// Find newest file in matches
	var newestTime time.Time
	var newestFile string
	for _, match := range matches {
		info, err := common.GetFileInfo(match)
		if err != nil {
			continue
		}
		if info.ModTime.After(newestTime) {
			newestTime = info.ModTime
			newestFile = match
		}
	}

	// Compare with source
	if sourceInfo.ModTime.After(newestTime) {
		relPath, _ := filepath.Rel(projectRoot, newestFile)
		return &Issue{
			Type:        "stale_build",
			Severity:    "error",
			Message:     fmt.Sprintf("%s is newer than build output (%s)", cmd.Source, relPath),
			FixAvailable: true,
			FixCommand:  getFixCommand(ecosystem, "stale_build"),
		}, nil
	}

	return nil, nil
}

// verifyCommand executes a command-based verification
func verifyCommand(cmd config.VerificationCommand, projectRoot string) (*Issue, error) {
	// TODO: Implement command execution verification
	// For now, return nil (no issue detected)
	return nil, nil
}

// getFixCommand retrieves the fix command for an issue type
func getFixCommand(ecosystem *detector.DetectedEcosystem, issueType string) string {
	cfg := ecosystem.Config
	for _, fix := range cfg.Ecosystem.Reconciliation.Fixes {
		if fix.IssueType == issueType {
			return fix.Command
		}
	}
	return ""
}

