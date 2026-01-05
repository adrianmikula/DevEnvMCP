package detector

import (
	"path/filepath"

	"dev-env-sentinel/internal/common"
	"dev-env-sentinel/internal/config"
)

// DetectedEcosystem represents a detected ecosystem in a project
type DetectedEcosystem struct {
	ID       string
	Config   *config.EcosystemConfig
	Confidence float64
	ProjectRoot string
}

// DetectEcosystems detects all ecosystems present in a project
func DetectEcosystems(projectRoot string, configs []*config.EcosystemConfig) ([]*DetectedEcosystem, error) {
	var detected []*DetectedEcosystem

	for _, cfg := range configs {
		if present, confidence := isEcosystemPresent(projectRoot, cfg); present {
			detected = append(detected, &DetectedEcosystem{
				ID:          cfg.Ecosystem.ID,
				Config:      cfg,
				Confidence:  confidence,
				ProjectRoot: projectRoot,
			})
		}
	}

	return detected, nil
}

// isEcosystemPresent checks if an ecosystem is present in a project
func isEcosystemPresent(projectRoot string, cfg *config.EcosystemConfig) (bool, float64) {
	detection := cfg.Ecosystem.Detection
	
	// Check required files
	requiredCount := 0
	for _, file := range detection.RequiredFiles {
		path := filepath.Join(projectRoot, file)
		if common.FileExists(path) {
			requiredCount++
		}
	}

	// All required files must be present
	if len(detection.RequiredFiles) > 0 && requiredCount < len(detection.RequiredFiles) {
		return false, 0
	}

	// Calculate confidence based on optional files and patterns
	confidence := 1.0
	if len(detection.RequiredFiles) > 0 {
		confidence = float64(requiredCount) / float64(len(detection.RequiredFiles))
	}

	// Boost confidence with optional files
	optionalCount := 0
	for _, file := range detection.OptionalFiles {
		path := filepath.Join(projectRoot, file)
		if common.FileExists(path) {
			optionalCount++
		}
	}
	if len(detection.OptionalFiles) > 0 {
		confidence += float64(optionalCount) / float64(len(detection.OptionalFiles)) * 0.2
		if confidence > 1.0 {
			confidence = 1.0
		}
	}

	// Check directory patterns
	patternCount := 0
	for _, pattern := range detection.DirectoryPatterns {
		expanded := common.ExpandPattern(pattern)
		path := filepath.Join(projectRoot, expanded)
		if common.DirExists(path) {
			patternCount++
		}
	}
	if len(detection.DirectoryPatterns) > 0 {
		confidence += float64(patternCount) / float64(len(detection.DirectoryPatterns)) * 0.1
		if confidence > 1.0 {
			confidence = 1.0
		}
	}

	return confidence >= 0.5, confidence
}

