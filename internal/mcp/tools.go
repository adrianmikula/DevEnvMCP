package mcp

import (
	"context"
	"fmt"

	"dev-env-sentinel/internal/auditor"
	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
	"dev-env-sentinel/internal/infra"
	"dev-env-sentinel/internal/reconciler"
	"dev-env-sentinel/internal/verifier"
)

// RegisterAllTools registers all MCP tools
func RegisterAllTools(server *Server, configs []*config.EcosystemConfig) {
	server.RegisterTool("verify_build_freshness", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleVerifyBuildFreshness(args, configs)
	})

	server.RegisterTool("check_infrastructure_parity", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleCheckInfrastructureParity(args, configs)
	})

	server.RegisterTool("env_var_audit", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleEnvVarAudit(args, configs)
	})

	server.RegisterTool("reconcile_environment", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleReconcileEnvironment(args, configs)
	})
}

// handleVerifyBuildFreshness handles the verify_build_freshness tool
func handleVerifyBuildFreshness(args map[string]interface{}, configs []*config.EcosystemConfig) (interface{}, error) {
	projectRoot, ok := args["project_root"].(string)
	if !ok {
		return nil, fmt.Errorf("project_root is required")
	}

	// Detect ecosystems
	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to detect ecosystems: %w", err)
	}

	if len(ecosystems) == 0 {
		return "No ecosystems detected in project", nil
	}

	// Verify build freshness for each ecosystem
	var reports []*verifier.FreshnessReport
	for _, eco := range ecosystems {
		report, err := verifier.VerifyBuildFreshness(projectRoot, eco)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	if len(reports) == 0 {
		return "No verification reports generated", nil
	}

	// Return first report (can be extended to return all)
	return reports[0], nil
}

// handleCheckInfrastructureParity handles the check_infrastructure_parity tool
func handleCheckInfrastructureParity(args map[string]interface{}, configs []*config.EcosystemConfig) (interface{}, error) {
	projectRoot, ok := args["project_root"].(string)
	if !ok {
		return nil, fmt.Errorf("project_root is required")
	}

	// Detect ecosystems
	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to detect ecosystems: %w", err)
	}

	if len(ecosystems) == 0 {
		return "No ecosystems detected in project", nil
	}

	// Check infrastructure for each ecosystem
	var reports []*infra.InfrastructureReport
	for _, eco := range ecosystems {
		report, err := infra.CheckInfrastructure(context.Background(), eco.Config)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	if len(reports) == 0 {
		return "No infrastructure reports generated", nil
	}

	// Return first report (can be extended to return all)
	return reports[0], nil
}

// handleEnvVarAudit handles the env_var_audit tool
func handleEnvVarAudit(args map[string]interface{}, configs []*config.EcosystemConfig) (interface{}, error) {
	projectRoot, ok := args["project_root"].(string)
	if !ok {
		return nil, fmt.Errorf("project_root is required")
	}

	// Detect ecosystems
	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to detect ecosystems: %w", err)
	}

	if len(ecosystems) == 0 {
		return "No ecosystems detected in project", nil
	}

	// Audit environment variables for each ecosystem
	var reports []*auditor.EnvVarReport
	for _, eco := range ecosystems {
		report, err := auditor.AuditEnvironmentVariables(projectRoot, eco.Config)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	if len(reports) == 0 {
		return "No environment variable reports generated", nil
	}

	// Return first report (can be extended to return all)
	return reports[0], nil
}

// handleReconcileEnvironment handles the reconcile_environment tool
func handleReconcileEnvironment(args map[string]interface{}, configs []*config.EcosystemConfig) (interface{}, error) {
	projectRoot, ok := args["project_root"].(string)
	if !ok {
		return nil, fmt.Errorf("project_root is required")
	}

	// Detect ecosystems
	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to detect ecosystems: %w", err)
	}

	if len(ecosystems) == 0 {
		return "No ecosystems detected in project", nil
	}

	// First, verify build freshness to get issues
	var allIssues []verifier.Issue
	for _, eco := range ecosystems {
		report, err := verifier.VerifyBuildFreshness(projectRoot, eco)
		if err != nil {
			continue
		}
		allIssues = append(allIssues, report.Issues...)
	}

	if len(allIssues) == 0 {
		return "No issues found to reconcile", nil
	}

	// Reconcile issues for first ecosystem (can be extended)
	report, err := reconciler.ReconcileEnvironment(context.Background(), projectRoot, allIssues, ecosystems[0])
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile environment: %w", err)
	}

	return report, nil
}

