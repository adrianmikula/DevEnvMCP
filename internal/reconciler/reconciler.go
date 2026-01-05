package reconciler

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
	"dev-env-sentinel/internal/verifier"
)

// ReconciliationReport contains reconciliation results
type ReconciliationReport struct {
	Fixed     []FixResult
	Failed    []FixResult
	IsSuccess bool
	Message   string
}

// FixResult represents the result of a fix attempt
type FixResult struct {
	IssueType string
	Command   string
	Success   bool
	Message   string
	Error     string
}

// ReconcileEnvironment reconciles environment issues
func ReconcileEnvironment(ctx context.Context, projectRoot string, issues []verifier.Issue, ecosystem *detector.DetectedEcosystem) (*ReconciliationReport, error) {
	report := &ReconciliationReport{
		Fixed:     []FixResult{},
		Failed:    []FixResult{},
		IsSuccess: true,
	}

	cfg := ecosystem.Config

	// Group issues by type and find fixes
	for _, issue := range issues {
		if !issue.FixAvailable {
			continue
		}

		fix := findFix(cfg, issue.Type)
		if fix == nil {
			report.Failed = append(report.Failed, FixResult{
				IssueType: issue.Type,
				Success:   false,
				Message:   "No fix available for this issue type",
			})
			report.IsSuccess = false
			continue
		}

		// Execute fix
		result := executeFix(ctx, projectRoot, fix, issue)
		if result.Success {
			report.Fixed = append(report.Fixed, result)
		} else {
			report.Failed = append(report.Failed, result)
			report.IsSuccess = false
		}
	}

	// Generate summary message
	if len(report.Fixed) > 0 {
		report.Message = fmt.Sprintf("Fixed %d issue(s)", len(report.Fixed))
	}
	if len(report.Failed) > 0 {
		if report.Message != "" {
			report.Message += ", "
		}
		report.Message += fmt.Sprintf("Failed to fix %d issue(s)", len(report.Failed))
	}

	return report, nil
}

// findFix finds a fix configuration for an issue type
func findFix(cfg *config.EcosystemConfig, issueType string) *config.Fix {
	for i := range cfg.Ecosystem.Reconciliation.Fixes {
		if cfg.Ecosystem.Reconciliation.Fixes[i].IssueType == issueType {
			return &cfg.Ecosystem.Reconciliation.Fixes[i]
		}
	}
	return nil
}

// executeFix executes a fix command
func executeFix(ctx context.Context, projectRoot string, fix *config.Fix, issue verifier.Issue) FixResult {
	result := FixResult{
		IssueType: fix.IssueType,
		Command:   fix.Command,
		Success:   false,
	}

	// Use fix command from config, or fall back to issue fix command
	command := fix.Command
	if command == "" {
		command = issue.FixCommand
	}

	if command == "" {
		result.Message = "No fix command available"
		return result
	}

	// Execute fix command
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Error = err.Error()
		result.Message = fmt.Sprintf("Fix command failed: %s", strings.TrimSpace(string(output)))
		return result
	}

	// Verify fix if verify command provided
	if fix.VerifyCommand != "" {
		verifyCtx, verifyCancel := context.WithTimeout(ctx, 1*time.Minute)
		defer verifyCancel()

		verifyCmd := exec.CommandContext(verifyCtx, "sh", "-c", fix.VerifyCommand)
		verifyCmd.Dir = projectRoot
		verifyOutput, verifyErr := verifyCmd.CombinedOutput()

		if verifyErr != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Fix executed but verification failed: %s", strings.TrimSpace(string(verifyOutput)))
			result.Error = verifyErr.Error()
			return result
		}

		result.Success = true
		result.Message = fmt.Sprintf("Fix executed and verified successfully: %s", fix.Description)
	} else {
		result.Success = true
		result.Message = fmt.Sprintf("Fix executed: %s", fix.Description)
	}

	return result
}

// ReconcileIssue reconciles a single issue
func ReconcileIssue(ctx context.Context, projectRoot string, issue verifier.Issue, ecosystem *detector.DetectedEcosystem) (*FixResult, error) {
	if !issue.FixAvailable {
		return nil, fmt.Errorf("no fix available for issue: %s", issue.Type)
	}

	cfg := ecosystem.Config
	fix := findFix(cfg, issue.Type)
	if fix == nil {
		return nil, fmt.Errorf("no fix configuration found for issue type: %s", issue.Type)
	}

	result := executeFix(ctx, projectRoot, fix, issue)
	return &result, nil
}

