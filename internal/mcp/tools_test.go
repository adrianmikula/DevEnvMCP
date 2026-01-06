package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"dev-env-sentinel/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// findConfigDir finds the ecosystem-configs directory relative to the project root
func findConfigDir() string {
	// Start from current directory and walk up to find ecosystem-configs
	dir, _ := os.Getwd()
	for {
		configDir := filepath.Join(dir, "ecosystem-configs")
		if info, err := os.Stat(configDir); err == nil && info.IsDir() {
			return configDir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root
		}
		dir = parent
	}
	return ""
}

func TestHandleVerifyBuildFreshness(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
	require.NoError(t, err)

	// Create src directory
	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	configDir := findConfigDir()
	if configDir == "" {
		t.Skip("ecosystem-configs directory not found")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	require.NoError(t, err)

	args := map[string]interface{}{
		"project_root": tmpDir,
	}

	result, err := handleVerifyBuildFreshness(args, configs)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleVerifyBuildFreshness_NoProjectRoot(t *testing.T) {
	configs := []*config.EcosystemConfig{}

	args := map[string]interface{}{
		// Missing project_root
	}

	_, err := handleVerifyBuildFreshness(args, configs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project_root is required")
}

func TestHandleVerifyBuildFreshness_NoEcosystems(t *testing.T) {
	tmpDir := t.TempDir()
	configs := []*config.EcosystemConfig{} // Empty configs

	args := map[string]interface{}{
		"project_root": tmpDir,
	}

	result, err := handleVerifyBuildFreshness(args, configs)
	require.NoError(t, err)
	assert.Equal(t, "No ecosystems detected in project", result)
}

func TestHandleCheckInfrastructureParity(t *testing.T) {
	tmpDir := t.TempDir()

	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	configDir := findConfigDir()
	if configDir == "" {
		t.Skip("ecosystem-configs directory not found")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	require.NoError(t, err)

	args := map[string]interface{}{
		"project_root": tmpDir,
	}

	result, err := handleCheckInfrastructureParity(args, configs)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleEnvVarAudit(t *testing.T) {
	tmpDir := t.TempDir()

	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	configDir := findConfigDir()
	if configDir == "" {
		t.Skip("ecosystem-configs directory not found")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	require.NoError(t, err)

	args := map[string]interface{}{
		"project_root": tmpDir,
	}

	result, err := handleEnvVarAudit(args, configs)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleReconcileEnvironment(t *testing.T) {
	tmpDir := t.TempDir()

	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	configDir := findConfigDir()
	if configDir == "" {
		t.Skip("ecosystem-configs directory not found")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	require.NoError(t, err)

	args := map[string]interface{}{
		"project_root": tmpDir,
	}

	server := NewServer()
	result, err := handleReconcileEnvironment(server, args, configs)
	require.NoError(t, err)
	
	// Should return "No issues found to reconcile" if no issues
	if str, ok := result.(string); ok {
		assert.Contains(t, str, "No issues")
	}
}

func TestRegisterAllTools(t *testing.T) {
	server := NewServer()
	configs := []*config.EcosystemConfig{}

	RegisterAllTools(server, configs)

	// Verify tools are registered
	assert.NotNil(t, server.tools["verify_build_freshness"])
	assert.NotNil(t, server.tools["check_infrastructure_parity"])
	assert.NotNil(t, server.tools["env_var_audit"])
	assert.NotNil(t, server.tools["reconcile_environment"])
}

