# Dev Environment Setup for Dev-Env Sentinel MCP

This document outlines the tools and libraries needed to build the Dev-Env Sentinel MCP server in Go.

## MCP Tools for Go

**⚠️ IMPORTANT**: See `docs/research/go-mcp-libraries.md` for comprehensive comparison of available Go MCP libraries.

### Top Candidates

#### 1. **Go MCP SDK** (Official - Preferred)
- **Status**: Official SDK from MCP maintainers
- **URL**: https://www.mcpstack.org/sdks/golang/go-mcp-sdk
- **Description**: Official Go SDK for building MCP servers
- **Features**:
  - Official support and protocol compliance
  - Type-safe APIs
  - High-performance concurrent server/client
  - Built-in middleware
  - Comprehensive testing utilities
  - Multiple transport layers (Stdio, HTTP, WebSocket, TCP)
- **Install**: Verify exact package name from repository
- **Decision**: **PREFERRED** if available and stable

#### 2. **mcp-golang** (Alternative)
- **URL**: https://github.com/metoro-io/mcp-golang
- **Description**: Community Go implementation with low boilerplate
- **Features**:
  - Type safety through native Go structs
  - Minimal boilerplate code
  - Modular design
  - Custom transports (stdio and HTTP)
  - Automatic schema generation
- **Install**: `go get github.com/metoro-io/mcp-golang`
- **Decision**: **ALTERNATIVE** if Go MCP SDK unavailable

#### 3. **mcp-go** (Fallback)
- **URL**: https://pkg.go.dev/gitlab.flexinfer.ai/libs/mcp-go
- **Description**: Production-proven MCP implementation
- **Features**:
  - Complete MCP protocol implementation
  - Full JSON-RPC 2.0 support
  - Multiple transports (Stdio, WebSocket)
  - Simple API
  - Production-ready
- **Install**: `go get gitlab.flexinfer.ai/libs/mcp-go`
- **Decision**: **FALLBACK** if others don't work

#### 4. **GoMCP** (Other Option)
- **URL**: https://gomcp.dev
- **Description**: Transport-agnostic MCP server library
- **Features**:
  - Transport flexibility (Stdio, SSE+HTTP, WebSocket, TCP)
  - Tool, resource, and prompt registration
  - Argument parsing utilities
- **Install**: Verify exact package name

### Decision Process

See `docs/research/mcp-library-decision.md` for:
- Decision tree
- Verification checklist
- Integration test approach
- Timeline

**Action Required**: Research and verify repositories before making final decision.

## Health Check Libraries

### 1. **health-go** (Recommended)
- **URL**: https://github.com/hellofresh/health-go
- **Description**: Comprehensive health check library with built-in checkers
- **Features**:
  - HTTP handler for health status
  - Generic checkers for:
    - RabbitMQ, PostgreSQL, Redis, HTTP
    - MongoDB, MySQL, gRPC
    - Memcached, InfluxDB, NATS
- **Install**: `go get github.com/hellofresh/health-go/v4`

### 2. **Checkup**
- **URL**: https://sourcegraph.github.io/checkup/
- **Description**: Self-hosted, distributed health check and status page system
- **Features**:
  - Cross-platform operation
  - Easily configurable
  - Can be embedded in any Go program
  - Supports HTTP, TCP, and DNS checks
- **Install**: `go get github.com/sourcegraph/checkup`

### 3. **gochecks**
- **URL**: https://github.com/aleasoluciones/gochecks
- **Description**: Utilities to check service health and publish events
- **Features**:
  - Scheduler for periodic checks
  - TCP port, ICMP/Ping, HTTP checks
  - SNMP get, RabbitMQ queue length, MySQL connectivity
  - Publishers (RabbitMQ/AMQP)
- **Install**: `go get github.com/aleasoluciones/gochecks`

## Fast Filesystem Scan Libraries

### 1. **fsnotify** (Recommended for file watching)
- **URL**: https://github.com/fsnotify/fsnotify
- **Description**: Cross-platform file system notifications
- **Features**:
  - Watch for file system changes
  - Cross-platform (Windows, Linux, macOS)
  - Event-driven filesystem monitoring
- **Install**: `go get github.com/fsnotify/fsnotify`

### 2. **Go Standard Library** (Recommended for fast scanning)
- **Package**: `filepath.WalkDir` (Go 1.16+)
- **Description**: Built-in efficient directory walking
- **Features**:
  - Faster than `filepath.Walk` (uses `os.ReadDir` internally)
  - No external dependencies
  - Platform-independent
- **Usage**: Part of `path/filepath` package (no install needed)

### 3. **fswalker** (For integrity checking)
- **URL**: https://google.github.io/fswalker/
- **Description**: File system integrity checking tool
- **Features**:
  - Walker component collects file system information
  - Reporter compares runs to detect differences
  - Useful for monitoring file system changes
- **Install**: `go get github.com/google/fswalker`

## Process Monitoring Libraries

### 1. **gopsutil** (Highly Recommended)
- **URL**: https://github.com/shirou/gopsutil
- **Description**: Cross-platform process and system utilization library
- **Features**:
  - Process information (PID, name, status, CPU, memory)
  - System metrics (CPU, memory, disk, network, sensors)
  - Platform-independent API
  - Similar to Python's psutil
- **Install**: `go get github.com/shirou/gopsutil/v3`
- **Usage**: Perfect for checking if processes are running, monitoring resource usage

### 2. **Go Standard Library** (Basic process operations)
- **Packages**: `os`, `os/exec`, `syscall`
- **Description**: Built-in process management
- **Features**:
  - Execute external commands
  - Process information via `os` package
  - Low-level system calls via `syscall`
- **Usage**: No install needed, part of standard library

## Recommended Setup

### Core Dependencies (Minimum Viable Setup)

```go
// MCP Framework (TBD - see docs/research/go-mcp-libraries.md)
// Option 1 (Preferred): github.com/modelcontextprotocol/go-sdk (verify)
// Option 2 (Alternative): github.com/metoro-io/mcp-golang
// Option 3 (Fallback): gitlab.flexinfer.ai/libs/mcp-go

// Health Checks
github.com/hellofresh/health-go/v4

// Process Monitoring
github.com/shirou/gopsutil/v3

// Filesystem (standard library is sufficient, but fsnotify for watching)
github.com/fsnotify/fsnotify

// Testing
github.com/stretchr/testify
```

**Note**: MCP library choice pending research verification. See `docs/research/mcp-library-decision.md` for decision process.

### Additional Utilities (As Needed)

```go
// For more advanced health checks
github.com/sourcegraph/checkup

// For filesystem integrity
github.com/google/fswalker
```

## Dependency Management

**✅ Go Modules is already set up!**

- **Status**: Go 1.24.4 installed
- **Module**: `go.mod` initialized as `dev-env-sentinel`
- **No installation needed**: Go Modules is built into Go

See `docs/research/go-dependency-management.md` for details.

### Adding Dependencies

```bash
# Add a dependency
go get <package-path>

# Example: Add testing library
go get github.com/stretchr/testify

# Clean up unused dependencies
go mod tidy
```

## Next Steps

1. **✅ Go Module Initialized**: Already done
2. **Install Core Dependencies**: Use `go get` commands (see below)
3. **Create Project Structure**: Set up the MCP server skeleton
4. **Implement Core Tools**:
   - `verify_build_freshness` (filesystem + timestamp checks)
   - `check_infrastructure_parity` (health checks + version verification)
   - `env_var_audit` (environment variable validation)
   - `reconcile_environment` (active mode - fix issues)

### Installing Dependencies

When ready, install dependencies with:

```bash
# Testing
go get github.com/stretchr/testify

# Process monitoring
go get github.com/shirou/gopsutil/v3

# Health checks
go get github.com/hellofresh/health-go/v4

# Filesystem watching
go get github.com/fsnotify/fsnotify

# MCP library (TBD - see docs/research/go-mcp-libraries.md)
# go get <mcp-library-package>

# Clean up
go mod tidy
```

## Notes

- **Go Standard Library**: Many operations (file walking, basic process info) can be done with standard library alone
- **Performance**: Go's `filepath.WalkDir` is highly optimized and often faster than third-party libraries
- **Cross-Platform**: All recommended libraries support Windows, Linux, and macOS
- **Static Binaries**: Go compiles to single static binaries - perfect for MCP server distribution

