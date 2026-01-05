package infra

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"dev-env-sentinel/internal/config"
)

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name      string
	Running   bool
	Version   string
	ExpectedVersion string
	Healthy   bool
	Message   string
}

// InfrastructureReport contains infrastructure check results
type InfrastructureReport struct {
	Services []ServiceStatus
	IsHealthy bool
	Issues   []string
}

// CheckInfrastructure checks infrastructure parity for an ecosystem
func CheckInfrastructure(ctx context.Context, cfg *config.EcosystemConfig) (*InfrastructureReport, error) {
	report := &InfrastructureReport{
		Services:  []ServiceStatus{},
		IsHealthy: true,
		Issues:    []string{},
	}

	for _, service := range cfg.Ecosystem.Infrastructure.Services {
		status, err := checkService(ctx, service)
		if err != nil {
			report.Issues = append(report.Issues, fmt.Sprintf("%s: %v", service.Name, err))
			continue
		}

		report.Services = append(report.Services, *status)

		if !status.Healthy {
			report.IsHealthy = false
			report.Issues = append(report.Issues, status.Message)
		}
	}

	return report, nil
}

// checkService checks a single service
func checkService(ctx context.Context, service config.Service) (*ServiceStatus, error) {
	status := &ServiceStatus{
		Name:    service.Name,
		Running: false,
		Healthy: false,
	}

	// Execute check command
	cmd := exec.CommandContext(ctx, "sh", "-c", service.CheckCommand)
	output, err := cmd.Output()
	if err != nil {
		status.Message = fmt.Sprintf("Service check failed: %v", err)
		return status, nil
	}

	status.Running = true
	outputStr := strings.TrimSpace(string(output))

	// Extract version if pattern provided
	if service.VersionExtract != "" {
		version, err := extractVersion(outputStr, service.VersionExtract)
		if err == nil {
			status.Version = version
		}
	}

	// If we got output, service is likely healthy
	if outputStr != "" {
		status.Healthy = true
		status.Message = fmt.Sprintf("%s is running", service.Name)
		if status.Version != "" {
			status.Message += fmt.Sprintf(" (version: %s)", status.Version)
		}
	} else {
		status.Message = fmt.Sprintf("%s check returned no output", service.Name)
	}

	return status, nil
}

// extractVersion extracts version from output using regex
func extractVersion(output, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("version pattern not found")
	}

	return matches[1], nil
}

// CheckServiceHealth checks if a service is healthy with timeout
func CheckServiceHealth(ctx context.Context, checkCommand string, timeout time.Duration) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", checkCommand)
	output, err := cmd.Output()
	if err != nil {
		return false, "", err
	}

	return true, strings.TrimSpace(string(output)), nil
}

