package auditor

import (
	"os"
	"path/filepath"
	"testing"

	"dev-env-sentinel/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditEnvironmentVariables(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a source file with environment variable references
	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	javaFile := filepath.Join(srcDir, "App.java")
	javaContent := `package com.example;
public class App {
    public static void main(String[] args) {
        String dbUrl = System.getenv("DATABASE_URL");
        String apiKey = System.getenv("API_KEY");
    }
}`
	err = os.WriteFile(javaFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "java-maven",
			Environment: config.Environment{
				VariablePatterns: []string{
					`System\.getenv\("([A-Z_][A-Z0-9_]*)"\)`,
				},
				ConfigFiles: []string{},
				RequiredVars: []string{},
			},
		},
	}

	// Set one environment variable
	os.Setenv("DATABASE_URL", "postgres://localhost/db")
	defer os.Unsetenv("DATABASE_URL")

	report, err := AuditEnvironmentVariables(tmpDir, cfg)
	require.NoError(t, err)
	require.NotNil(t, report)

	// Should find 2 references
	assert.Len(t, report.References, 2)
	
	// DATABASE_URL should be set, API_KEY should be missing
	foundDbUrl := false
	foundApiKey := false
	for _, ref := range report.References {
		if ref.Name == "DATABASE_URL" {
			foundDbUrl = true
			assert.True(t, ref.IsSet)
			assert.Equal(t, "postgres://localhost/db", ref.Value)
		}
		if ref.Name == "API_KEY" {
			foundApiKey = true
			assert.False(t, ref.IsSet)
		}
	}
	assert.True(t, foundDbUrl)
	assert.True(t, foundApiKey)

	// Should not be healthy (missing API_KEY)
	assert.False(t, report.IsHealthy)
	assert.Contains(t, report.Missing, "API_KEY")
}

func TestAuditEnvironmentVariables_AllSet(t *testing.T) {
	tmpDir := t.TempDir()

	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	javaFile := filepath.Join(srcDir, "App.java")
	javaContent := `package com.example;
public class App {
    public static void main(String[] args) {
        String dbUrl = System.getenv("DATABASE_URL");
    }
}`
	err = os.WriteFile(javaFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "java-maven",
			Environment: config.Environment{
				VariablePatterns: []string{
					`System\.getenv\("([A-Z_][A-Z0-9_]*)"\)`,
				},
			},
		},
	}

	os.Setenv("DATABASE_URL", "postgres://localhost/db")
	defer os.Unsetenv("DATABASE_URL")

	report, err := AuditEnvironmentVariables(tmpDir, cfg)
	require.NoError(t, err)

	assert.True(t, report.IsHealthy)
	assert.Empty(t, report.Missing)
}

func TestAuditEnvironmentVariables_ConfigFileVars(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .env file
	envFile := filepath.Join(tmpDir, ".env")
	envContent := `DATABASE_URL=postgres://localhost/db
API_KEY=secret123
# Comment
OTHER_VAR=value`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	cfg := &config.EcosystemConfig{
		Ecosystem: config.Ecosystem{
			ID: "test",
			Environment: config.Environment{
				VariablePatterns: []string{},
				ConfigFiles:      []string{".env"},
				RequiredVars:      []string{},
			},
		},
	}

	// Set DATABASE_URL but not API_KEY
	os.Setenv("DATABASE_URL", "postgres://localhost/db")
	defer os.Unsetenv("DATABASE_URL")

	report, err := AuditEnvironmentVariables(tmpDir, cfg)
	require.NoError(t, err)

	// Should detect missing API_KEY and OTHER_VAR from config file
	assert.False(t, report.IsHealthy)
	assert.Contains(t, report.Missing, "API_KEY")
	assert.Contains(t, report.Missing, "OTHER_VAR")
}

func TestFindEnvVarReferences(t *testing.T) {
	tmpDir := t.TempDir()

	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	goFile := filepath.Join(srcDir, "main.go")
	goContent := `package main
import "os"
func main() {
    dbUrl := os.Getenv("DATABASE_URL")
    apiKey := os.Getenv("API_KEY")
}`
	err = os.WriteFile(goFile, []byte(goContent), 0644)
	require.NoError(t, err)

	patterns := []string{
		`os\.Getenv\("([A-Z_][A-Z0-9_]*)"\)`,
	}

	refs, err := findEnvVarReferences(tmpDir, patterns)
	require.NoError(t, err)
	assert.Len(t, refs, 2)

	// Verify references
	names := make(map[string]bool)
	for _, ref := range refs {
		names[ref.Name] = true
		assert.Contains(t, ref.File, "main.go")
		assert.Greater(t, ref.Line, 0)
	}
	assert.True(t, names["DATABASE_URL"])
	assert.True(t, names["API_KEY"])
}

func TestFindPatternMatches(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		pattern string
		want    []string
	}{
		{
			name:    "single match",
			line:    `String url = System.getenv("DATABASE_URL");`,
			pattern: `System\.getenv\("([A-Z_][A-Z0-9_]*)"\)`,
			want:    []string{"DATABASE_URL"},
		},
		{
			name:    "multiple matches",
			line:    `String url = System.getenv("DB_URL"); String key = System.getenv("API_KEY");`,
			pattern: `System\.getenv\("([A-Z_][A-Z0-9_]*)"\)`,
			want:    []string{"DB_URL", "API_KEY"},
		},
		{
			name:    "no matches",
			line:    `String url = "hardcoded";`,
			pattern: `System\.getenv\("([A-Z_][A-Z0-9_]*)"\)`,
			want:    nil, // Function returns nil when no matches
		},
		{
			name:    "invalid pattern",
			line:    `String url = System.getenv("DB_URL");`,
			pattern: `[invalid regex`,
			want:    nil, // Function returns nil for invalid patterns
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := findPatternMatches(tt.line, tt.pattern)
			assert.Equal(t, tt.want, matches)
		})
	}
}

func TestIsSourceFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"file.go", true},
		{"file.java", true},
		{"file.js", true},
		{"file.ts", true},
		{"file.py", true},
		{"file.cpp", true},
		{"file.c", true},
		{"file.h", true},
		{"file.cs", true},
		{"file.txt", false},
		{"file.yaml", false},
		{"file", false},
		{"FILE.GO", true}, // Case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.expected, isSourceFile(tt.path))
		})
	}
}

func TestParseConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Test .env file parsing
	envFile := filepath.Join(tmpDir, ".env")
	envContent := `DATABASE_URL=postgres://localhost/db
API_KEY=secret123
# Comment line
EMPTY_VAR=
OTHER_VAR=value with spaces`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	vars, err := parseConfigFile(envFile)
	require.NoError(t, err)
	assert.Len(t, vars, 4)
	assert.Contains(t, vars, "DATABASE_URL")
	assert.Contains(t, vars, "API_KEY")
	assert.Contains(t, vars, "EMPTY_VAR")
	assert.Contains(t, vars, "OTHER_VAR")
}

func TestParseConfigFile_NonEnvFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Test non-.env file (should return empty)
	otherFile := filepath.Join(tmpDir, "config.properties")
	err := os.WriteFile(otherFile, []byte("key=value"), 0644)
	require.NoError(t, err)

	vars, err := parseConfigFile(otherFile)
	require.NoError(t, err)
	assert.Empty(t, vars)
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{[]string{"a"}, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			assert.Equal(t, tt.expected, contains(tt.slice, tt.item))
		})
	}
}

func TestFindEnvVarReferences_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create node_modules directory (should be skipped)
	nodeModules := filepath.Join(tmpDir, "node_modules", "package")
	err := os.MkdirAll(nodeModules, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(nodeModules, "file.js"), []byte(`process.env.API_KEY`), 0644)
	require.NoError(t, err)

	// Create source file outside node_modules
	srcFile := filepath.Join(tmpDir, "src", "main.go")
	err = os.MkdirAll(filepath.Dir(srcFile), 0755)
	require.NoError(t, err)
	err = os.WriteFile(srcFile, []byte(`os.Getenv("DATABASE_URL")`), 0644)
	require.NoError(t, err)

	patterns := []string{`os\.Getenv\("([A-Z_][A-Z0-9_]*)"\)`}

	refs, err := findEnvVarReferences(tmpDir, patterns)
	require.NoError(t, err)

	// Should only find DATABASE_URL, not API_KEY from node_modules
	assert.Len(t, refs, 1)
	assert.Equal(t, "DATABASE_URL", refs[0].Name)
}

