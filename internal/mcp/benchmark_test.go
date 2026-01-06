//go:build benchmark
// +build benchmark

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

// setupLargeProjectBench creates a large project structure for benchmarking
func setupLargeProjectBench(b *testing.B, numFiles int, depth int) string {
	tmpDir := b.TempDir()

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

	// Create source directory structure
	srcDir := filepath.Join(tmpDir, "src", "main", "java", "com", "example")
	os.MkdirAll(srcDir, 0755)

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
		createNestedDirsBench(b, tmpDir, depth, 10)
	}

	return tmpDir
}

// createNestedDirsBench creates nested directory structures for benchmarks
func createNestedDirsBench(b *testing.B, baseDir string, depth, filesPerDir int) {
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
		createNestedDirsBench(b, dir, depth-1, filesPerDir)
	}
}

// BenchmarkBootstrapTime benchmarks the MCP server bootstrap time
func BenchmarkBootstrapTime(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
			b.Fatal(err)
		}

		// Create server
		server := NewServer()
		RegisterAllTools(server, configs)

		elapsed := time.Since(start)
		b.Logf("Bootstrap time: %v", elapsed)

		// Ensure bootstrap is fast (< 100ms)
		if elapsed > 100*time.Millisecond {
			b.Logf("WARNING: Bootstrap took %v, should be < 100ms", elapsed)
		}
	}
}

// BenchmarkEcosystemDetection_Small benchmarks ecosystem detection on small projects
func BenchmarkEcosystemDetection_Small(b *testing.B) {
	projectRoot := setupLargeProjectBench(b, 10, 0)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectEcosystems(projectRoot, configs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEcosystemDetection_Medium benchmarks ecosystem detection on medium projects
func BenchmarkEcosystemDetection_Medium(b *testing.B) {
	projectRoot := setupLargeProjectBench(b, 100, 2)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectEcosystems(projectRoot, configs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEcosystemDetection_Large benchmarks ecosystem detection on large projects
func BenchmarkEcosystemDetection_Large(b *testing.B) {
	projectRoot := setupLargeProjectBench(b, 1000, 3)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectEcosystems(projectRoot, configs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuildFreshness_Small benchmarks build freshness verification on small projects
func BenchmarkBuildFreshness_Small(b *testing.B) {
	projectRoot := setupLargeProjectBench(b, 10, 0)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		b.Fatal(err)
	}

	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	if err != nil {
		b.Fatal(err)
	}
	if len(ecosystems) == 0 {
		b.Skip("No ecosystems detected")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := verifier.VerifyBuildFreshness(projectRoot, ecosystems[0])
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuildFreshness_Large benchmarks build freshness verification on large projects
func BenchmarkBuildFreshness_Large(b *testing.B) {
	projectRoot := setupLargeProjectBench(b, 1000, 3)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		b.Fatal(err)
	}

	ecosystems, err := detector.DetectEcosystems(projectRoot, configs)
	if err != nil {
		b.Fatal(err)
	}
	if len(ecosystems) == 0 {
		b.Skip("No ecosystems detected")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := verifier.VerifyBuildFreshness(projectRoot, ecosystems[0])
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMCPTool_VerifyBuildFreshness benchmarks the MCP tool execution
func BenchmarkMCPTool_VerifyBuildFreshness(b *testing.B) {
	projectRoot := setupLargeProjectBench(b, 100, 2)
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}
	configs, err := config.DiscoverEcosystemConfigs(configDir)
	if err != nil {
		b.Fatal(err)
	}

	args := map[string]interface{}{
		"project_root": projectRoot,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handleVerifyBuildFreshness(args, configs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConfigLoading benchmarks config file loading
func BenchmarkConfigLoading(b *testing.B) {
	configDir := filepath.Join("..", "..", "ecosystem-configs")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = filepath.Join("..", "ecosystem-configs")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := config.DiscoverEcosystemConfigs(configDir)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestBootstrapTime_Requirement_Benchmark tests that bootstrap completes in acceptable time (benchmark version)
func TestBootstrapTime_Requirement_Benchmark(t *testing.T) {
	// This test is duplicated in performance_test.go
	// Skip in benchmark mode to avoid conflicts
	t.Skip("Use BenchmarkBootstrapTime for benchmarks")
}

// TestToolExecutionTime_Requirement_Benchmark tests that tool execution is fast (benchmark version)
func TestToolExecutionTime_Requirement_Benchmark(t *testing.T) {
	// This test is duplicated in performance_test.go
	// Skip in benchmark mode to avoid conflicts
	t.Skip("Use BenchmarkMCPTool_VerifyBuildFreshness for benchmarks")
}

// TestLargeProjectPerformance_Benchmark tests performance on very large projects (benchmark version)
func TestLargeProjectPerformance_Benchmark(t *testing.T) {
	// This test is duplicated in performance_test.go
	// Skip in benchmark mode to avoid conflicts
	t.Skip("Use BenchmarkEcosystemDetection_Large and BenchmarkBuildFreshness_Large for benchmarks")
}

