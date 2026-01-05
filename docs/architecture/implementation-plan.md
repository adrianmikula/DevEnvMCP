# Implementation Plan

This document outlines how to implement the language-agnostic, configuration-driven architecture.

## Phase 1: Core Engine (Language-Agnostic)

### 1.1 Configuration System

**Components**:
- `config/loader.go` - Load and parse YAML configs
- `config/schema.go` - Go structs matching YAML schema
- `config/validator.go` - Validate config files
- `config/merger.go` - Merge default/user/project configs

**Key Functions**:
```go
// Load ecosystem configuration
func LoadEcosystemConfig(ecosystemID string) (*EcosystemConfig, error)

// Find all ecosystem configs in directory
func DiscoverEcosystemConfigs(configDir string) ([]*EcosystemConfig, error)

// Merge configs (default -> user -> project)
func MergeConfigs(configs ...*EcosystemConfig) *EcosystemConfig
```

### 1.2 Ecosystem Detector

**Components**:
- `detector/detector.go` - Main detection logic
- `detector/pattern_matcher.go` - Match files/patterns from config

**Key Functions**:
```go
// Detect ecosystems in a project directory
func DetectEcosystems(projectRoot string, configs []*EcosystemConfig) ([]*DetectedEcosystem, error)

// Check if ecosystem is present based on config
func IsEcosystemPresent(projectRoot string, config *EcosystemConfig) (bool, error)
```

**Generic Logic**:
1. Scan project directory for manifest files
2. For each config, check required/optional files
3. Score confidence based on matches
4. Return detected ecosystems with confidence scores

### 1.3 Build Freshness Verifier

**Components**:
- `verifier/freshness.go` - Build freshness verification
- `verifier/timestamp.go` - Timestamp comparison utilities
- `verifier/command_executor.go` - Execute verification commands

**Key Functions**:
```go
// Verify build freshness for an ecosystem
func VerifyBuildFreshness(projectRoot string, ecosystem *DetectedEcosystem) (*FreshnessReport, error)

// Compare timestamps (generic)
func CompareTimestamps(source, target string) (bool, error) // true if source is newer

// Execute verification command from config
func ExecuteVerificationCommand(cmd *VerificationCommand, projectRoot string) (*CommandResult, error)
```

**Generic Logic**:
1. Load ecosystem config
2. For each verification command in config:
   - Execute command based on type (timestamp_compare, command, etc.)
   - Collect results
3. Aggregate results into report

### 1.4 Dependency Auditor

**Components**:
- `auditor/dependency.go` - Dependency auditing
- `auditor/parser.go` - Parse dependency information (generic)

**Key Functions**:
```go
// Audit dependencies for an ecosystem
func AuditDependencies(projectRoot string, ecosystem *DetectedEcosystem) (*DependencyReport, error)

// Parse manifest file (format-agnostic)
func ParseManifest(manifestPath string, format string) (*Manifest, error)
```

### 1.5 Environment Variable Auditor

**Components**:
- `auditor/envvar.go` - Environment variable auditing
- `auditor/code_scanner.go` - Scan code for env var references

**Key Functions**:
```go
// Audit environment variables
func AuditEnvironmentVariables(projectRoot string, ecosystem *DetectedEcosystem) (*EnvVarReport, error)

// Find env var references in code
func FindEnvVarReferences(projectRoot string, patterns []string) ([]*EnvVarReference, error)

// Check if variables are set
func CheckVariablesSet(references []*EnvVarReference) (*EnvVarStatus, error)
```

### 1.6 Infrastructure Parity Checker

**Components**:
- `infra/checker.go` - Infrastructure checking
- `infra/service.go` - Service health/version checking

**Key Functions**:
```go
// Check infrastructure parity
func CheckInfrastructure(ecosystem *DetectedEcosystem) (*InfrastructureReport, error)

// Check service health
func CheckService(service *ServiceConfig) (*ServiceStatus, error)

// Get service version
func GetServiceVersion(service *ServiceConfig) (string, error)
```

### 1.7 Environment Reconciler

**Components**:
- `reconciler/reconciler.go` - Environment reconciliation
- `reconciler/fix_executor.go` - Execute fix commands

**Key Functions**:
```go
// Reconcile environment issues
func ReconcileEnvironment(projectRoot string, issues []*Issue, ecosystem *DetectedEcosystem) (*ReconciliationReport, error)

// Execute fix command
func ExecuteFix(fix *FixConfig, projectRoot string) (*FixResult, error)
```

## Phase 2: MCP Integration

### 2.1 MCP Server Setup

**Components**:
- `mcp/server.go` - MCP server implementation
- `mcp/tools.go` - MCP tool definitions

**MCP Tools**:
1. `verify_build_freshness` - Check build artifact freshness
2. `check_infrastructure_parity` - Verify services/versions
3. `env_var_audit` - Audit environment variables
4. `reconcile_environment` - Auto-fix issues

### 2.2 Tool Implementations

Each MCP tool:
1. Detects ecosystems in project
2. Loads ecosystem configs
3. Executes generic verification logic
4. Returns formatted results

## Phase 3: Configuration Files

### 3.1 Default Configurations

Create YAML configs for:
- `java-maven.yaml`
- `java-gradle.yaml`
- `npm.yaml`
- `yarn.yaml`
- `pnpm.yaml`

### 3.2 Configuration Validation

- Schema validation on load
- Runtime validation of commands
- Error messages for invalid configs

## Phase 4: Testing Strategy

### 4.1 Unit Tests

- Test each building block independently
- Mock configs for testing
- Test config loading/merging

### 4.2 Integration Tests

- Test with real project directories
- Test ecosystem detection
- Test verification commands

### 4.3 Configuration Tests

- Validate all default configs
- Test config schema compliance
- Test config merging logic

## File Structure

```
dev-env-sentinel/
├── cmd/
│   └── sentinel/
│       └── main.go              # MCP server entry point
├── internal/
│   ├── config/
│   │   ├── loader.go
│   │   ├── schema.go
│   │   ├── validator.go
│   │   └── merger.go
│   ├── detector/
│   │   ├── detector.go
│   │   └── pattern_matcher.go
│   ├── verifier/
│   │   ├── freshness.go
│   │   ├── timestamp.go
│   │   └── command_executor.go
│   ├── auditor/
│   │   ├── dependency.go
│   │   ├── envvar.go
│   │   └── code_scanner.go
│   ├── infra/
│   │   ├── checker.go
│   │   └── service.go
│   ├── reconciler/
│   │   ├── reconciler.go
│   │   └── fix_executor.go
│   └── mcp/
│       ├── server.go
│       └── tools.go
├── ecosystem-configs/
│   ├── java-maven.yaml
│   ├── java-gradle.yaml
│   ├── npm.yaml
│   ├── yarn.yaml
│   └── pnpm.yaml
├── docs/
│   └── architecture/
├── go.mod
└── go.sum
```

## Key Design Decisions

### 1. Command Execution
- Use `os/exec` for executing commands
- Support both shell commands and direct binary execution
- Capture stdout/stderr for parsing
- Timeout support for long-running commands

### 2. Pattern Matching
- Use Go's `filepath.Match` for glob patterns
- Use `regexp` for regex patterns
- Support both absolute and relative paths
- Handle platform-specific path separators

### 3. Timestamp Comparison
- Use `os.Stat` to get file modification times
- Compare timestamps with configurable precision
- Handle directory timestamps (use newest file)

### 4. Error Handling
- Structured errors with context
- Continue verification even if one check fails
- Aggregate errors for reporting

### 5. Configuration Caching
- Cache loaded configs in memory
- Invalidate cache on config file changes (optional)
- Support hot-reloading configs (optional)

## Extension Points

### Adding a New Ecosystem

1. Create YAML config file in `ecosystem-configs/`
2. Follow schema defined in `configuration-schema.md`
3. Test with real project
4. No code changes required!

### Adding a New Verification Type

1. Add command type to schema
2. Implement handler in `verifier/command_executor.go`
3. Update schema documentation
4. Test with example config

### Custom Verification Logic

Users can:
1. Create custom config in `.sentinel/configs/`
2. Override default verification commands
3. Add project-specific checks
4. Extend existing ecosystem configs

