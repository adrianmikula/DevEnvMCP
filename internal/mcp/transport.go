package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Transport defines the interface for MCP transport layers
type Transport interface {
	Start(ctx context.Context, server *Server) error
}

// StdioTransport implements stdio-based transport (for local use)
type StdioTransport struct{}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport() *StdioTransport {
	return &StdioTransport{}
}

// Start starts the server with stdio transport
func (t *StdioTransport) Start(ctx context.Context, server *Server) error {
	// Initialize MCP protocol
	if err := server.initialize(); err != nil {
		return err
	}

	// Start message loop
	return server.messageLoop()
}

// SSETransport implements SSE+HTTP transport (for Apify/cloud deployments)
type SSETransport struct {
	port     string
	readOnly bool // If true, only handles reads (for SSE)
}

// NewSSETransport creates a new SSE transport
func NewSSETransport(port string) *SSETransport {
	return &SSETransport{
		port:     port,
		readOnly: false,
	}
}

// Start starts the server with SSE+HTTP transport
func (t *SSETransport) Start(ctx context.Context, server *Server) error {
	// Set up HTTP handlers
	http.HandleFunc("/sse", t.handleSSE(server))
	http.HandleFunc("/message", t.handleMessage(server))
	http.HandleFunc("/health", t.handleHealth)

	addr := ":" + t.port
	if t.port == "" {
		addr = ":8080" // Default port
	}

	fmt.Fprintf(os.Stderr, "Starting MCP server with SSE transport on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleSSE handles Server-Sent Events connections
func (t *SSETransport) handleSSE(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Flush headers
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		// Send initial connection message
		fmt.Fprintf(w, "data: %s\n\n", `{"type":"connected"}`)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		// Keep connection alive (Apify will handle the actual message flow)
		<-r.Context().Done()
	}
}

// handleMessage handles HTTP POST messages (for sending requests to server)
func (t *SSETransport) handleMessage(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Read request body
		var msg map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Handle the message
		var response map[string]interface{}
		if method, ok := msg["method"].(string); ok {
			switch method {
			case "initialize":
				// Handle initialize
				initResp := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      msg["id"],
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
				response = initResp
			case "tools/list":
				response = server.handleToolsListResponse(msg)
			case "tools/call":
				response = server.handleToolCallResponse(msg)
			default:
				response = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      msg["id"],
					"error": map[string]interface{}{
						"code":    -32601,
						"message": fmt.Sprintf("Method not found: %s", method),
					},
				}
			}
		} else {
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      msg["id"],
				"error": map[string]interface{}{
					"code":    -32600,
					"message": "Invalid Request",
				},
			}
		}

		// Send response
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(response)
	}
}

// handleHealth handles health check requests
func (t *SSETransport) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","transport":"sse"}`)
}

// handleToolsListResponse handles tools/list and returns response map
func (s *Server) handleToolsListResponse(msg map[string]interface{}) map[string]interface{} {
	tools := []map[string]interface{}{}

	for name := range s.tools {
		tools = append(tools, map[string]interface{}{
			"name":        name,
			"description": getToolDescription(name),
		})
	}

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      msg["id"],
		"result": map[string]interface{}{
			"tools": tools,
		},
	}
}

// handleToolCallResponse handles tools/call and returns response map
func (s *Server) handleToolCallResponse(msg map[string]interface{}) map[string]interface{} {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      msg["id"],
			"error": map[string]interface{}{
				"code":    -32602,
				"message": "Invalid params",
			},
		}
	}

	name, ok := params["name"].(string)
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      msg["id"],
			"error": map[string]interface{}{
				"code":    -32602,
				"message": "Invalid tool name",
			},
		}
	}

	handler, ok := s.tools[name]
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      msg["id"],
			"error": map[string]interface{}{
				"code":    -32601,
				"message": fmt.Sprintf("Unknown tool: %s", name),
			},
		}
	}

	args, _ := params["arguments"].(map[string]interface{})

	// Execute tool
	result, err := handler(context.Background(), args)
	if err != nil {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      msg["id"],
			"error": map[string]interface{}{
				"code":    -1,
				"message": err.Error(),
			},
		}
	}

	return map[string]interface{}{
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
}

// DetectTransport detects which transport to use based on environment
func DetectTransport() Transport {
	// Check for SSE/HTTP mode (for Apify)
	if port := os.Getenv("SENTINEL_HTTP_PORT"); port != "" {
		return NewSSETransport(port)
	}

	// Check for explicit transport
	if transport := os.Getenv("SENTINEL_TRANSPORT"); transport == "sse" || transport == "http" {
		port := os.Getenv("SENTINEL_HTTP_PORT")
		if port == "" {
			port = "8080"
		}
		return NewSSETransport(port)
	}

	// Default to stdio
	return NewStdioTransport()
}

