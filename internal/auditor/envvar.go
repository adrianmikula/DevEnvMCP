package auditor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"dev-env-sentinel/internal/common"
	"dev-env-sentinel/internal/config"
)

// EnvVarReference represents a reference to an environment variable
type EnvVarReference struct {
	Name      string
	File      string
	Line      int
	Pattern   string
	IsSet     bool
	Value     string
}

// EnvVarReport contains environment variable audit results
type EnvVarReport struct {
	References []EnvVarReference
	Missing    []string
	IsHealthy  bool
	Issues     []string
}

// AuditEnvironmentVariables audits environment variables for an ecosystem
func AuditEnvironmentVariables(projectRoot string, cfg *config.EcosystemConfig) (*EnvVarReport, error) {
	report := &EnvVarReport{
		References: []EnvVarReference{},
		Missing:    []string{},
		IsHealthy:  true,
		Issues:     []string{},
	}

	// Find all environment variable references in code
	refs, err := findEnvVarReferences(projectRoot, cfg.Ecosystem.Environment.VariablePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to find env var references: %w", err)
	}

	// Check which variables are set
	missing := make(map[string]bool)
	for i := range refs {
		ref := &refs[i]
		if value, exists := os.LookupEnv(ref.Name); exists {
			ref.IsSet = true
			ref.Value = value
		} else {
			ref.IsSet = false
			missing[ref.Name] = true
			report.IsHealthy = false
		}
		report.References = append(report.References, *ref)
	}

	// Convert missing map to slice
	for name := range missing {
		report.Missing = append(report.Missing, name)
		report.Issues = append(report.Issues, fmt.Sprintf("Missing environment variable: %s", name))
	}

	// Check config files for declared variables
	configVars, err := findConfigFileVars(projectRoot, cfg.Ecosystem.Environment.ConfigFiles)
	if err == nil {
		for _, varName := range configVars {
			if _, exists := os.LookupEnv(varName); !exists {
				if !contains(report.Missing, varName) {
					report.Missing = append(report.Missing, varName)
					report.Issues = append(report.Issues, fmt.Sprintf("Variable %s declared in config but not set", varName))
					report.IsHealthy = false
				}
			}
		}
	}

	return report, nil
}

// findEnvVarReferences finds environment variable references in code
func findEnvVarReferences(projectRoot string, patterns []string) ([]EnvVarReference, error) {
	var refs []EnvVarReference

	// Walk through source directories
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip non-source files
		if info.IsDir() {
			// Skip common non-source directories
			if strings.Contains(path, "node_modules") || 
			   strings.Contains(path, ".git") ||
			   strings.Contains(path, "target") ||
			   strings.Contains(path, "build") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check source files
		if !isSourceFile(path) {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Check each pattern
		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			for _, pattern := range patterns {
				matches := findPatternMatches(line, pattern)
				for _, match := range matches {
					refs = append(refs, EnvVarReference{
						Name:    match,
						File:    path,
						Line:    lineNum + 1,
						Pattern: pattern,
						IsSet:   false,
					})
				}
			}
		}

		return nil
	})

	return refs, err
}

// findPatternMatches finds matches for a regex pattern in a line
func findPatternMatches(line, pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	matches := re.FindAllStringSubmatch(line, -1)
	var result []string
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1]) // First capture group
		}
	}

	return result
}

// isSourceFile checks if a file is a source file
func isSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	sourceExts := []string{".go", ".java", ".js", ".ts", ".jsx", ".tsx", ".py", ".cpp", ".c", ".h", ".cs"}
	for _, se := range sourceExts {
		if ext == se {
			return true
		}
	}
	return false
}

// findConfigFileVars finds variables declared in config files
func findConfigFileVars(projectRoot string, configFiles []string) ([]string, error) {
	var vars []string

	for _, pattern := range configFiles {
		expanded := common.ExpandPattern(pattern)
		fullPattern := filepath.Join(projectRoot, expanded)

		matches, err := common.FindFilesByPattern(fullPattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			fileVars, err := parseConfigFile(match)
			if err != nil {
				continue
			}
			vars = append(vars, fileVars...)
		}
	}

	return vars, nil
}

// parseConfigFile parses a config file for environment variables
func parseConfigFile(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var vars []string
	lines := strings.Split(string(content), "\n")

	// Simple parsing for .env files (KEY=VALUE format)
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".env" || strings.Contains(path, ".env") {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			if idx := strings.Index(line, "="); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				if key != "" {
					vars = append(vars, key)
				}
			}
		}
	}

	return vars, nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

