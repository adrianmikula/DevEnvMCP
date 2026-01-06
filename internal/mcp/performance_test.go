package mcp

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/detector"
	"dev-env-sentinel/internal/verifier"
)

// setupLargeProject creates a large project structure for performance testing
func setupLargeProject(t *testing.T, numFiles int, depth int) string {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>large-project</artifactId>
    <version>1.0.0</version>
</project>`
	os.WriteFile(pomPath, []byte(pomContent), 0644)

	// Create source directory structure (required for detection)
	srcDir := filepath.Join(tmpDir, "src", "main", "java", "com", "example")
	os.MkdirAll(srcDir, 0755)
	
	// Also create test directory (optional but helps with detection)
	testDir := filepath.Join(tmpDir, "src", "test", "java", "com", "example")
	os.MkdirAll(testDir, 0755)

	// Create many Java files
	for i := 0; i < numFiles; i++ {
		javaFile := filepath.Join(srcDir, "Class"+string(rune(65+i%26))+string(rune(48+i/26))+".java")
		content := `package com.example; public class Class` + string(rune(65+i%26)) + string(rune(48+i/26)) + ` { }`
		os.WriteFile(javaFile, []byte(content), 0644)
	}

	// Create build output directory with many class files
	targetDir := filepath.Join(tmpDir, "target", "classes", "com", "example")
	os.MkdirAll(targetDir, 0755)
	for i := 0; i < numFiles; i++ {
		classFile := filepath.Join(targetDir, "Class"+string(rune(65+i%26))+string(rune(48+i/26))+".class")
		os.WriteFile(classFile, []byte("fake class file"), 0644)
	}

	// Create nested directories if depth > 0
	if depth > 0 {
		createNestedDirs(t, tmpDir, depth, 10)
	}

	return tmpDir
}

// createNestedDirs creates nested directory structures
func createNestedDirs(t *testing.T, baseDir string, depth, filesPerDir int) {
	if depth <= 0 {
		return
	}

	for i := 0; i < filesPerDir; i++ {
		dir := filepath.Join(baseDir, "dir"+string(rune(48+i)))
		os.MkdirAll(dir, 0755)

		// Create some files in this directory
		for j := 0; j < filesPerDir; j++ {
			file := filepath.Join(dir, "file"+string(rune(48+j))+".txt")
			os.WriteFile(file, []byte("test content"), 0644)
		}

		// Recurse
		createNestedDirs(t, dir, depth-1, filesPerDir)
	}
}

// TestBootstrapTime_Requirement tests that bootstrap completes in acceptable time
func TestBootstrapTime_Requirement(t *testing.T) {
	start := time.Now()

	// Load configs (try multiple paths)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			configDir = filepath.Join(wd, "..", "..", "ecosystem-configs")
		}
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create server
	server := NewServer()
	RegisterAllTools(server, configs)

	elapsed := time.Since(start)

	// Bootstrap should complete in < 100ms
	maxTime := 100 * time.Millisecond
	if elapsed > maxTime {
		t.Errorf("Bootstrap took %v, should be < %v", elapsed, maxTime)
	} else {
		t.Logf("✓ Bootstrap completed in %v (requirement: < %v)", elapsed, maxTime)
	}
}

// TestToolExecutionTime_Requirement tests that tool execution is fast
func TestToolExecutionTime_Requirement(t *testing.T) {
	projectRoot := setupLargeProject(t, 100, 2)
	// Load configs (try multiple paths)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			configDir = filepath.Join(wd, "..", "..", "ecosystem-configs")
		}
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		t.Fatal(err)
	}

	args := map[string]interface{}{
		"project_root": projectRoot,
	}

	// Test verify_build_freshness tool
	start := time.Now()
	_, err = handleVerifyBuildFreshness(args, configs)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	// Tool execution should complete in < 1 second for medium projects
	maxTime := 1 * time.Second
	if elapsed > maxTime {
		t.Errorf("Tool execution took %v, should be < %v", elapsed, maxTime)
	} else {
		t.Logf("✓ Tool execution completed in %v (requirement: < %v)", elapsed, maxTime)
	}
}

// TestLargeProjectPerformance tests performance on very large projects
func TestLargeProjectPerformance(t *testing.T) {
	// Create a large project (1000 files, 3 levels deep)
	projectRoot := setupLargeProject(t, 1000, 3)
	// Load configs (try multiple paths)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			configDir = filepath.Join(wd, "..", "..", "ecosystem-configs")
		}
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		t.Fatal(err)
	}

	// Test ecosystem detection
	start := time.Now()
	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	detectionTime := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	if len(ecosystems) == 0 {
		t.Skip("No ecosystems detected")
	}

	// Test build freshness verification
	start = time.Now()
	_, err = verifier.VerifyBuildFreshness(projectRoot, ecosystems[0])
	verificationTime := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	// Large projects should still complete in reasonable time (< 2 seconds)
	maxDetectionTime := 2 * time.Second
	maxVerificationTime := 2 * time.Second

	if detectionTime > maxDetectionTime {
		t.Errorf("Ecosystem detection took %v, should be < %v", detectionTime, maxDetectionTime)
	} else {
		t.Logf("✓ Ecosystem detection: %v (requirement: < %v)", detectionTime, maxDetectionTime)
	}

	if verificationTime > maxVerificationTime {
		t.Errorf("Build freshness verification took %v, should be < %v", verificationTime, maxVerificationTime)
	} else {
		t.Logf("✓ Build freshness verification: %v (requirement: < %v)", verificationTime, maxVerificationTime)
	}
}

// TestEcosystemDetectionPerformance_Small tests detection on small projects
func TestEcosystemDetectionPerformance_Small(t *testing.T) {
	projectRoot := setupLargeProject(t, 10, 0)
	// Load configs (try multiple paths)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			configDir = filepath.Join(wd, "..", "..", "ecosystem-configs")
		}
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	_, err = detector.DetectEcosystems(projectRoot, configs)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	// Small projects should be very fast (< 100ms)
	maxTime := 100 * time.Millisecond
	if elapsed > maxTime {
		t.Errorf("Small project detection took %v, should be < %v", elapsed, maxTime)
	} else {
		t.Logf("✓ Small project detection: %v (requirement: < %v)", elapsed, maxTime)
	}
}

// TestEcosystemDetectionPerformance_Medium tests detection on medium projects
func TestEcosystemDetectionPerformance_Medium(t *testing.T) {
	projectRoot := setupLargeProject(t, 100, 2)
	// Load configs (try multiple paths)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			configDir = filepath.Join(wd, "..", "..", "ecosystem-configs")
		}
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	_, err = detector.DetectEcosystems(projectRoot, configs)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	// Medium projects should complete in < 500ms
	maxTime := 500 * time.Millisecond
	if elapsed > maxTime {
		t.Errorf("Medium project detection took %v, should be < %v", elapsed, maxTime)
	} else {
		t.Logf("✓ Medium project detection: %v (requirement: < %v)", elapsed, maxTime)
	}
}

