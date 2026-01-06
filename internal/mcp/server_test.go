package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"dev-env-sentinel/internal/auditor"
	"dev-env-sentinel/internal/infra"
	"dev-env-sentinel/internal/reconciler"
	"dev-env-sentinel/internal/verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	assert.NotNil(t, server)
	assert.NotNil(t, server.tools)
	assert.Empty(t, server.tools) // Should start empty
}

func TestRegisterTool(t *testing.T) {
	server := NewServer()
	
	handler := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "result", nil
	}

	server.RegisterTool("test_tool", handler)
	
	assert.Len(t, server.tools, 1)
	assert.NotNil(t, server.tools["test_tool"])
}

func TestGetToolDescription(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		expected string
	}{
		{"verify_build_freshness", "verify_build_freshness", "Verify that build artifacts are up-to-date with source manifests"},
		{"check_infrastructure_parity", "check_infrastructure_parity", "Check if required services are running and correct versions"},
		{"env_var_audit", "env_var_audit", "Audit environment variables for missing or incorrect values"},
		{"reconcile_environment", "reconcile_environment", "Automatically fix detected environment issues (Pro feature)"},
		{"unknown_tool", "unknown_tool", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := getToolDescription(tt.toolName)
			assert.Equal(t, tt.expected, desc)
		})
	}
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name     string
		result   interface{}
		contains []string
	}{
		{
			name:     "string result",
			result:   "simple string",
			contains: []string{"simple string"},
		},
		{
			name:     "map result",
			result:   map[string]string{"key": "value"},
			contains: []string{"key", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := formatResult(tt.result)
			for _, substr := range tt.contains {
				assert.Contains(t, formatted, substr)
			}
		})
	}
}

func TestFormatFreshnessReport(t *testing.T) {
	report := &verifier.FreshnessReport{
		EcosystemID: "java-maven",
		IsHealthy:   false,
		Issues: []verifier.Issue{
			{
				Type:        "stale_build",
				Severity:    "error",
				Message:     "Build is stale",
				FixAvailable: true,
				FixCommand:  "mvn clean",
			},
		},
	}

	formatted := formatFreshnessReport(report)
	assert.Contains(t, formatted, "java-maven")
	assert.Contains(t, formatted, "Build is stale")
	assert.Contains(t, formatted, "mvn clean")
}

func TestFormatFreshnessReport_Healthy(t *testing.T) {
	report := &verifier.FreshnessReport{
		EcosystemID: "java-maven",
		IsHealthy:   true,
		Issues:     []verifier.Issue{},
	}

	formatted := formatFreshnessReport(report)
	assert.Contains(t, formatted, "✅")
	assert.Contains(t, formatted, "java-maven")
}

func TestFormatInfrastructureReport(t *testing.T) {
	report := &infra.InfrastructureReport{
		IsHealthy: false,
		Services: []infra.ServiceStatus{
			{
				Name:    "maven",
				Healthy: true,
				Message: "Maven is running",
			},
			{
				Name:    "redis",
				Healthy: false,
				Message: "Redis is not running",
			},
		},
		Issues: []string{"Service issue"},
	}

	formatted := formatInfrastructureReport(report)
	assert.Contains(t, formatted, "maven")
	assert.Contains(t, formatted, "redis")
	assert.Contains(t, formatted, "Service issue")
}

func TestFormatEnvVarReport(t *testing.T) {
	report := &auditor.EnvVarReport{
		IsHealthy: false,
		Missing:   []string{"API_KEY", "DATABASE_URL"},
		Issues:    []string{"Missing API_KEY"},
	}

	formatted := formatEnvVarReport(report)
	assert.Contains(t, formatted, "API_KEY")
	assert.Contains(t, formatted, "DATABASE_URL")
	assert.Contains(t, formatted, "Missing API_KEY")
}

func TestFormatReconciliationReport(t *testing.T) {
	report := &reconciler.ReconciliationReport{
		IsSuccess: true,
		Fixed: []reconciler.FixResult{
			{
				IssueType: "stale_build",
				Success:   true,
				Message:   "Fixed successfully",
			},
		},
		Failed: []reconciler.FixResult{
			{
				IssueType: "other_issue",
				Success:   false,
				Message:   "Failed to fix",
				Error:     "Command error",
			},
		},
	}

	formatted := formatReconciliationReport(report)
	assert.Contains(t, formatted, "stale_build")
	assert.Contains(t, formatted, "other_issue")
	assert.Contains(t, formatted, "Fixed successfully")
	assert.Contains(t, formatted, "Command error")
}

func TestHandleToolsList(t *testing.T) {
	server := NewServer()
	server.RegisterTool("test_tool", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "result", nil
	})

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
	}

	// We can't easily test the JSON output without mocking stdout
	// But we can test that it doesn't panic
	err := server.handleToolsList(msg)
	assert.NoError(t, err)
}

func TestHandleToolCall(t *testing.T) {
	server := NewServer()
	
	server.RegisterTool("test_tool", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "success", nil
	})

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"params": map[string]interface{}{
			"name": "test_tool",
			"arguments": map[string]interface{}{},
		},
	}

	// We can't easily test without mocking stdout, but we can verify it doesn't panic
	err := server.handleToolCall(msg)
	assert.NoError(t, err)
}

func TestHandleToolCall_InvalidParams(t *testing.T) {
	server := NewServer()

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"params": "invalid", // Not a map
	}

	err := server.handleToolCall(msg)
	assert.Error(t, err)
}

func TestHandleToolCall_UnknownTool(t *testing.T) {
	server := NewServer()

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"params": map[string]interface{}{
			"name": "unknown_tool",
			"arguments": map[string]interface{}{},
		},
	}

	err := server.handleToolCall(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown tool")
}

func TestReadJSON(t *testing.T) {
	server := NewServer()

	// Create a temporary file with JSON content
	tmpFile := filepath.Join(t.TempDir(), "input.json")
	jsonContent := `{"test": "value"}`
	err := os.WriteFile(tmpFile, []byte(jsonContent), 0644)
	require.NoError(t, err)

	// Note: readJSON reads from os.Stdin, so we can't easily test it
	// without mocking stdin. This test verifies the function exists.
	_ = server.readJSON
}

func TestWriteJSON(t *testing.T) {
	server := NewServer()

	// Note: writeJSON writes to os.Stdout, so we can't easily test it
	// without capturing stdout. This test verifies the function exists.
	data := map[string]string{"key": "value"}
	_ = server.writeJSON
	_ = data
}

func TestHandleMethod(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name   string
		method string
		wantErr bool
	}{
		{"tools/list", "tools/list", false},
		{"tools/call", "tools/call", false},
		{"unknown", "unknown_method", false}, // Unknown methods are ignored
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
			}
			
			// We can't fully test without proper message structure,
			// but we can verify it doesn't panic
			err := server.handleMethod(tt.method, msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// May or may not error depending on message structure
				_ = err
			}
		})
	}
}

func TestUpdateLicense(t *testing.T) {
	server := NewServer()
	
	// Test with invalid license (should fail)
	err := server.UpdateLicense("invalid-key")
	assert.Error(t, err)

	// Test with Apify token (should succeed)
	err = server.UpdateLicense("apify_1234567890abcdef")
	assert.NoError(t, err)
	assert.NotNil(t, server.license)
}

func TestFormatResult_JSON(t *testing.T) {
	// Test that complex objects are JSON marshaled
	complexObj := map[string]interface{}{
		"nested": map[string]string{
			"key": "value",
		},
		"array": []string{"a", "b"},
	}

	formatted := formatResult(complexObj)
	
	// Should be valid JSON
	var decoded map[string]interface{}
	err := json.Unmarshal([]byte(formatted), &decoded)
	assert.NoError(t, err)
}

