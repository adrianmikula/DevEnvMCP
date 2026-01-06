package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dev-env-sentinel/internal/config"
	"dev-env-sentinel/internal/mcp"
)

func main() {
	// Check if running as MCP server (no args) or CLI mode
	if len(os.Args) == 1 {
		// MCP server mode
		runMCPServer()
	} else {
		// CLI mode (for testing)
		runCLIMode()
	}
}

// getConfigBaseDir returns the base directory for config discovery
func getConfigBaseDir() string {
	// Check for explicit config directory in environment
	if configDir := os.Getenv("SENTINEL_CONFIG_DIR"); configDir != "" {
		return configDir
	}

	// Try to find config relative to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// Check if config directory exists relative to executable
		configPath := filepath.Join(exeDir, "config")
		if _, err := os.Stat(configPath); err == nil {
			return exeDir
		}
		// For npm package, config might be in parent directory
		parentConfigPath := filepath.Join(exeDir, "..", "config")
		if _, err := os.Stat(parentConfigPath); err == nil {
			return filepath.Join(exeDir, "..")
		}
	}

	// Fallback to current working directory
	return "."
}

// runMCPServer runs the MCP server
func runMCPServer() {
	// Get base directory for config discovery
	baseDir := getConfigBaseDir()
	
	// Load ecosystem configs from config directory structure
	configs, err := config.DiscoverEcosystemConfigs(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading configs from %s: %v\n", baseDir, err)
		os.Exit(1)
	}

	// Create MCP server
	server := mcp.NewServer()

	// Register all tools
	mcp.RegisterAllTools(server, configs)

	// Start server
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v\n", err)
		os.Exit(1)
	}
}

// runCLIMode runs in CLI mode for testing
func runCLIMode() {
	fmt.Fprintf(os.Stderr, "CLI mode not yet implemented\n")
	fmt.Fprintf(os.Stderr, "Run without arguments to start MCP server\n")
	os.Exit(1)
}

