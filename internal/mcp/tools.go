package mcp

import (
	"context"
	"fmt"

	"dev-env-sentinel/internal/auditor"
	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
	"dev-env-sentinel/internal/infra"
	"dev-env-sentinel/internal/license"
	"dev-env-sentinel/internal/reconciler"
	"dev-env-sentinel/internal/verifier"
)

// RegisterAllTools registers all MCP tools
func RegisterAllTools(server *Server, configs []*config.EcosystemConfig) {
	// Free tier tools
	server.RegisterTool("verify_build_freshness", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleVerifyBuildFreshness(args, configs)
	})

	server.RegisterTool("check_infrastructure_parity", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleCheckInfrastructureParity(args, configs)
	})

	server.RegisterTool("env_var_audit", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleEnvVarAudit(args, configs)
	})

	// Premium tier tool (gated)
	server.RegisterTool("reconcile_environment", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleReconcileEnvironment(server, args, configs)
	})

	// Monetization tools
	server.RegisterTool("get_pro_license", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleGetProLicense(server)
	})

	server.RegisterTool("activate_pro", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleActivatePro(server, args)
	})

	server.RegisterTool("check_license_status", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return handleCheckLicenseStatus(server)
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

// handleReconcileEnvironment handles the reconcile_environment tool (PREMIUM FEATURE)
func handleReconcileEnvironment(server *Server, args map[string]interface{}, configs []*config.EcosystemConfig) (interface{}, error) {
	// Check if feature is available
	if err := server.featureManager.RequireFeature("reconcile_environment"); err != nil {
		upgradeMsg := server.featureManager.GetUpgradeMessage("reconcile_environment")
		return upgradeMsg, fmt.Errorf("premium feature not available: %w", err)
	}

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

// handleGetProLicense returns information about getting a Pro license
func handleGetProLicense(server *Server) (interface{}, error) {
	stripeLink := license.GetStripePaymentLink()
	apifyURL := license.GetApifyActorURL()
	
	msg := fmt.Sprintf(
		"🚀 Upgrade to Dev-Env Sentinel Pro\n\n"+
			"Unlock powerful features:\n"+
			"• Auto-fix environment issues\n"+
			"• Advanced diagnostics\n"+
			"• Docker orchestration (Enterprise)\n"+
			"• Priority support\n\n"+
			"Purchase Options:\n\n"+
			"1. Stripe Payment Link (One-time/Subscription):\n   %s\n\n"+
			"2. Apify Actor (Pay-Per-Event):\n   %s\n\n"+
			"After purchasing, use the 'activate_pro' tool with your license key.",
		stripeLink, apifyURL,
	)
	
	return msg, nil
}

// handleActivatePro activates a Pro license
func handleActivatePro(server *Server, args map[string]interface{}) (interface{}, error) {
	key, ok := args["license_key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("license_key is required")
	}

	// Update license on server
	if err := server.UpdateLicense(key); err != nil {
		return nil, fmt.Errorf("failed to activate license: %w", err)
	}

	// Get updated license info
	lic := server.license
	msg := fmt.Sprintf(
		"✅ License activated successfully!\n\n"+
			"Tier: %s\n"+
			"Status: %s\n"+
			"Features enabled: %d\n\n"+
			"You now have access to all Pro features, including auto-fix capabilities.",
		lic.Tier,
		map[bool]string{true: "Valid", false: "Invalid"}[lic.IsValid],
		len(lic.Features),
	)

	return msg, nil
}

// handleCheckLicenseStatus returns current license status
func handleCheckLicenseStatus(server *Server) (interface{}, error) {
	lic := server.license
	
	status := "Free"
	if lic.IsValid {
		status = fmt.Sprintf("%s (Valid)", lic.Tier)
		if lic.ExpiresAt != nil {
			status += fmt.Sprintf(" - Expires: %s", lic.ExpiresAt.Format("2006-01-02"))
		} else if lic.Tier != "free" {
			status += " - Lifetime"
		}
	} else if lic.Tier != "free" {
		status = fmt.Sprintf("%s (Invalid/Expired)", lic.Tier)
	}

	msg := fmt.Sprintf(
		"License Status: %s\n\n"+
			"Available Features:\n",
		status,
	)

	for _, feature := range lic.Features {
		msg += fmt.Sprintf("• %s\n", feature)
	}

	if !lic.IsValid && lic.Tier != "free" {
		msg += "\n⚠️ Your license is invalid or expired. Use 'get_pro_license' to purchase a new one."
	} else if lic.Tier == "free" {
		msg += "\n💡 Upgrade to Pro to unlock auto-fix and advanced features. Use 'get_pro_license' for details."
	}

	return msg, nil
}

