package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"dev-env-sentinel/internal/verifier"
)

// Server represents the MCP server
type Server struct {
	tools map[string]ToolHandler
}

// ToolHandler is a function that handles a tool call
type ToolHandler func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// NewServer creates a new MCP server
func NewServer() *Server {
	return &Server{
		tools: make(map[string]ToolHandler),
	}
}

// RegisterTool registers a tool handler
func (s *Server) RegisterTool(name string, handler ToolHandler) {
	s.tools[name] = handler
}

// Start starts the MCP server with stdio transport
func (s *Server) Start() error {
	// Initialize MCP protocol
	if err := s.initialize(); err != nil {
		return err
	}

	// Start message loop
	return s.messageLoop()
}

// initialize sends the initialize request/response
func (s *Server) initialize() error {
	// Read initialize request
	var initReq map[string]interface{}
	if err := s.readJSON(&initReq); err != nil {
		return err
	}

	// Send initialize response
	initResp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      initReq["id"],
		"result": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "dev-env-sentinel",
				"version": "0.1.0",
			},
		},
	}

	return s.writeJSON(initResp)
}

// messageLoop processes incoming messages
func (s *Server) messageLoop() error {
	for {
		var msg map[string]interface{}
		if err := s.readJSON(&msg); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// Handle different message types
		if method, ok := msg["method"].(string); ok {
			if err := s.handleMethod(method, msg); err != nil {
				// Log error but continue
				continue
			}
		}
	}
}

// handleMethod handles a method call
func (s *Server) handleMethod(method string, msg map[string]interface{}) error {
	switch method {
	case "tools/list":
		return s.handleToolsList(msg)
	case "tools/call":
		return s.handleToolCall(msg)
	default:
		// Unknown method - ignore
		return nil
	}
}

// handleToolsList handles the tools/list request
func (s *Server) handleToolsList(msg map[string]interface{}) error {
	tools := []map[string]interface{}{}

	for name := range s.tools {
		tools = append(tools, map[string]interface{}{
			"name":        name,
			"description": getToolDescription(name),
		})
	}

	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      msg["id"],
		"result": map[string]interface{}{
			"tools": tools,
		},
	}

	return s.writeJSON(resp)
}

// handleToolCall handles a tool call request
func (s *Server) handleToolCall(msg map[string]interface{}) error {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid params")
	}

	name, ok := params["name"].(string)
	if !ok {
		return fmt.Errorf("invalid tool name")
	}

	handler, ok := s.tools[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}

	args, _ := params["arguments"].(map[string]interface{})

	// Execute tool
	result, err := handler(context.Background(), args)
	if err != nil {
		// Send error response
		resp := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      msg["id"],
			"error": map[string]interface{}{
				"code":    -1,
				"message": err.Error(),
			},
		}
		return s.writeJSON(resp)
	}

	// Send success response
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      msg["id"],
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": formatResult(result),
				},
			},
		},
	}

	return s.writeJSON(resp)
}

// readJSON reads a JSON message from stdin
func (s *Server) readJSON(v interface{}) error {
	decoder := json.NewDecoder(os.Stdin)
	return decoder.Decode(v)
}

// writeJSON writes a JSON message to stdout
func (s *Server) writeJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// getToolDescription returns the description for a tool
func getToolDescription(name string) string {
	descriptions := map[string]string{
		"verify_build_freshness":    "Verify that build artifacts are up-to-date with source manifests",
		"check_infrastructure_parity": "Check if required services are running and correct versions",
		"env_var_audit":            "Audit environment variables for missing or incorrect values",
		"reconcile_environment":     "Automatically fix detected environment issues",
	}
	return descriptions[name]
}

// formatResult formats a result for MCP response
func formatResult(result interface{}) string {
	switch v := result.(type) {
	case string:
		return v
	case *verifier.FreshnessReport:
		return formatFreshnessReport(v)
	default:
		data, _ := json.MarshalIndent(v, "", "  ")
		return string(data)
	}
}

// formatFreshnessReport formats a freshness report
func formatFreshnessReport(report *verifier.FreshnessReport) string {
	if report.IsHealthy {
		return fmt.Sprintf("✅ Build freshness check passed for %s", report.EcosystemID)
	}

	msg := fmt.Sprintf("❌ Build freshness issues found for %s:\n\n", report.EcosystemID)
	for _, issue := range report.Issues {
		msg += fmt.Sprintf("- %s: %s\n", issue.Severity, issue.Message)
		if issue.FixAvailable {
			msg += fmt.Sprintf("  Fix: %s\n", issue.FixCommand)
		}
	}
	return msg
}

