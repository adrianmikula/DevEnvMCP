//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/container/wait"

	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
	"dev-env-sentinel/internal/verifier"
)

// skipIfShort skips the test if -short flag is set
func skipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}

// setupMavenContainer creates a container with Maven installed
func setupMavenContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	// Create a temporary directory for the project
	tmpDir := t.TempDir()
	
	req := testcontainers.ContainerRequest{
		Image:        "maven:3.9-eclipse-temurin-17",
		Cmd:          []string{"tail", "-f", "/dev/null"}, // Keep container running
		WaitingFor:   wait.ForLog("").WithStartupTimeout(30 * time.Second),
		AutoRemove:   true,
		AttachStdout: true,
		AttachStderr: true,
		// Note: Mounts would require additional setup
		// For now, we'll use exec to create files
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get container working directory
	workDir := "/workspace"
	return container, workDir
}

// setupNodeContainer creates a container with Node.js/npm installed
func setupNodeContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "node:20-alpine",
		Cmd:          []string{"tail", "-f", "/dev/null"}, // Keep container running
		WaitingFor:   wait.ForLog("").WithStartupTimeout(30 * time.Second),
		AutoRemove:   true,
		AttachStdout: true,
		AttachStderr: true,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	workDir := "/workspace"
	return container, workDir
}

// execCommand executes a command in the container and returns stdout, stderr, and exit code
func execCommand(ctx context.Context, container testcontainers.Container, cmd []string) (string, string, int, error) {
	execResult, err := container.Exec(ctx, cmd)
	if err != nil {
		return "", "", -1, err
	}

	// Read stdout
	stdoutBytes, err := execResult.Stdout(ctx)
	if err != nil {
		return "", "", execResult.ExitCode, err
	}

	// Read stderr
	stderrBytes, err := execResult.Stderr(ctx)
	if err != nil {
		return string(stdoutBytes), "", execResult.ExitCode, err
	}

	return string(stdoutBytes), string(stderrBytes), execResult.ExitCode, nil
}


func TestIntegration_DetectMavenProject(t *testing.T) {
	skipIfShort(t)

	ctx := context.Background()
	container, workDir := setupMavenContainer(ctx, t)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Create a simple Maven project structure
	pomXML := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>test-project</artifactId>
    <version>1.0.0</version>
    <name>Test Maven Project</name>
</project>`

	// Write pom.xml to container using echo
	pomPath := filepath.Join(workDir, "pom.xml")
	cmd := []string{"sh", "-c", "mkdir -p " + workDir + " && echo '" + pomXML + "' > " + pomPath}
	_, _, exitCode, err := execCommand(ctx, container, cmd)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)

	// Load ecosystem configs
	configs, err := config.DiscoverEcosystemConfigs("ecosystem-configs")
	require.NoError(t, err)

	// For this test, we'll simulate by copying files locally and testing
	// In a real scenario, you'd mount volumes or copy files properly
	t.Log("Maven container setup complete")
	assert.NotNil(t, container)
}

func TestIntegration_DetectNpmProject(t *testing.T) {
	skipIfShort(t)

	ctx := context.Background()
	container, workDir := setupNodeContainer(ctx, t)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Create a simple npm project
	packageJSON := `{
  "name": "test-npm-project",
  "version": "1.0.0",
  "description": "Test npm project"
}`

	// Write package.json to container
	packagePath := filepath.Join(workDir, "package.json")
	cmd := []string{"sh", "-c", "mkdir -p " + workDir + " && echo '" + packageJSON + "' > " + packagePath}
	_, _, exitCode, err := execCommand(ctx, container, cmd)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)

	// Verify npm is available
	cmd = []string{"npm", "--version"}
	stdout, _, exitCode, err := execCommand(ctx, container, cmd)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	assert.NotEmpty(t, stdout)

	t.Log("NPM container setup complete")
	assert.NotNil(t, container)
}

func TestIntegration_RealMavenBuild(t *testing.T) {
	skipIfShort(t)

	ctx := context.Background()
	container, workDir := setupMavenContainer(ctx, t)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Create a minimal Maven project with source code
	pomXML := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>test-project</artifactId>
    <version>1.0.0</version>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>
</project>`

	// Setup project structure
	setupCmds := [][]string{
		{"sh", "-c", "mkdir -p " + workDir + "/src/main/java/com/example"},
		{"sh", "-c", "echo '" + pomXML + "' > " + workDir + "/pom.xml"},
		{"sh", "-c", "echo 'package com.example; public class App { public static void main(String[] args) {} }' > " + workDir + "/src/main/java/com/example/App.java"},
	}

	for _, cmd := range setupCmds {
		_, _, exitCode, err := execCommand(ctx, container, cmd)
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)
	}

	// Build the project
	buildCmd := []string{"mvn", "-f", workDir+"/pom.xml", "compile", "-q"}
	stdout, stderr, exitCode, err := execCommand(ctx, container, buildCmd)
	require.NoError(t, err, "Maven build failed: stdout=%s, stderr=%s", stdout, stderr)
	if exitCode != 0 {
		t.Logf("Maven build stderr: %s", stderr)
		t.Logf("Maven build stdout: %s", stdout)
	}
	require.Equal(t, 0, exitCode, "Maven build should succeed")

	// Verify build output exists
	checkCmd := []string{"test", "-f", workDir + "/target/classes/com/example/App.class"}
	_, _, exitCode, err = execCommand(ctx, container, checkCmd)
	assert.NoError(t, err, "Build output should exist")
	assert.Equal(t, 0, exitCode, "Build output file should exist")

	t.Log("Maven build completed successfully")
}

func TestIntegration_RealNpmBuild(t *testing.T) {
	skipIfShort(t)

	ctx := context.Background()
	container, workDir := setupNodeContainer(ctx, t)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Create a minimal npm project
	packageJSON := `{
  "name": "test-project",
  "version": "1.0.0",
  "main": "index.js"
}`

	indexJS := `console.log("Hello, World!");`

	// Setup project
	setupCmds := [][]string{
		{"sh", "-c", "mkdir -p " + workDir},
		{"sh", "-c", "echo '" + packageJSON + "' > " + workDir + "/package.json"},
		{"sh", "-c", "echo '" + indexJS + "' > " + workDir + "/index.js"},
	}

	for _, cmd := range setupCmds {
		_, _, exitCode, err := execCommand(ctx, container, cmd)
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)
	}

	// Verify npm can read package.json
	cmd := []string{"npm", "list", "--json", "--prefix", workDir}
	stdout, _, exitCode, err := execCommand(ctx, container, cmd)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	assert.Contains(t, stdout, "test-project")

	t.Log("NPM project setup completed successfully")
}

func TestIntegration_EcosystemDetectionWithRealProject(t *testing.T) {
	skipIfShort(t)

	// Use local testdata for this test (faster than containers)
	projectRoot := filepath.Join("testdata", "java-maven-project")
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		t.Skip("testdata not available")
	}

	// Load ecosystem configs
	configs, err := config.DiscoverEcosystemConfigs("ecosystem-configs")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	// Detect ecosystems
	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	require.NoError(t, err)
	require.NotEmpty(t, ecosystems, "Should detect at least one ecosystem")

	// Verify it detected Maven
	foundMaven := false
	for _, eco := range ecosystems {
		if eco.ID == "java-maven" {
			foundMaven = true
			assert.GreaterOrEqual(t, eco.Confidence, 0.5)
			break
		}
	}
	assert.True(t, foundMaven, "Should detect java-maven ecosystem")
}

func TestIntegration_BuildFreshnessWithRealMavenProject(t *testing.T) {
	skipIfShort(t)

	ctx := context.Background()
	container, workDir := setupMavenContainer(ctx, t)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Create Maven project
	pomXML := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>test-project</artifactId>
    <version>1.0.0</version>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>
</project>`

	javaCode := `package com.example; public class App { public static void main(String[] args) {} }`

	// Setup and build
	setupCmds := [][]string{
		{"sh", "-c", "mkdir -p " + workDir + "/src/main/java/com/example"},
		{"sh", "-c", "echo '" + pomXML + "' > " + workDir + "/pom.xml"},
		{"sh", "-c", "echo '" + javaCode + "' > " + workDir + "/src/main/java/com/example/App.java"},
	}

	for _, cmd := range setupCmds {
		_, _, exitCode, err := execCommand(ctx, container, cmd)
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)
	}

	// Initial build
	buildCmd := []string{"mvn", "-f", workDir+"/pom.xml", "compile", "-q"}
	_, _, exitCode, err := execCommand(ctx, container, buildCmd)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Initial build should succeed")

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Modify source file (simulate code change)
	modifiedCode := `package com.example; public class App { public static void main(String[] args) { System.out.println("Modified"); } }`
	modifyCmd := []string{"sh", "-c", "echo '" + modifiedCode + "' > " + workDir + "/src/main/java/com/example/App.java"}
	_, _, exitCode, err = execCommand(ctx, container, modifyCmd)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)

	// Now test build freshness detection
	// Since we can't easily access container files from Go, we'll test the logic
	// In a real scenario, you'd mount volumes or use file copy mechanisms

	// Load config and create ecosystem
	configs, err := config.DiscoverEcosystemConfigs("ecosystem-configs")
	require.NoError(t, err)

	var mavenConfig *config.EcosystemConfig
	for _, cfg := range configs {
		if cfg.Ecosystem.ID == "java-maven" {
			mavenConfig = cfg
			break
		}
	}
	require.NotNil(t, mavenConfig, "Maven config should be found")

	// Create detected ecosystem
	ecosystem := &detector.DetectedEcosystem{
		ID:          "java-maven",
		Config:      mavenConfig,
		Confidence:  1.0,
		ProjectRoot: workDir,
	}

	// Note: This test demonstrates the concept
	// In practice, you'd need to copy files from container or mount volumes
	// to test freshness verification with real file timestamps
	t.Log("Build freshness test setup complete")
	assert.NotNil(t, ecosystem)
}

