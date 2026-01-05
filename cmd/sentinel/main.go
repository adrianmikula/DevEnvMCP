package main

import (
	"fmt"
	"os"

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

// runMCPServer runs the MCP server
func runMCPServer() {
	// Load ecosystem configs
	configs, err := config.DiscoverEcosystemConfigs("ecosystem-configs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading configs: %v\n", err)
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

