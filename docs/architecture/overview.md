# Dev-Env Sentinel MCP Architecture

## Design Philosophy

The Dev-Env Sentinel MCP is built on a **language-agnostic, configuration-driven architecture**. All language and tool-specific logic is externalized into YAML configuration files, allowing the core engine to remain generic and extensible.

## Core Principles

1. **Separation of Concerns**: Core engine handles generic dev environment concepts; language/tool specifics live in config
2. **Extensibility**: Adding support for new languages/tools requires only adding a YAML config file
3. **Maintainability**: Tool-specific changes don't require code changes
4. **Testability**: Each building block can be tested independently

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    MCP Server (Go)                           │
│  ┌───────────────────────────────────────────────────────┐  │
│  │           Core Engine (Language-Agnostic)              │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐ │  │
│  │  │   Detector   │  │   Verifier   │  │  Reconciler │ │  │
│  │  └──────────────┘  └──────────────┘  └─────────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
│                          │                                   │
│                          ▼                                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         Configuration Loader & Parser                 │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │   YAML Configuration Files           │
        │  - ecosystem-configs/                │
        │    ├── java-maven.yaml              │
        │    ├── java-gradle.yaml             │
        │    ├── npm.yaml                     │
        │    ├── yarn.yaml                    │
        │    └── pnpm.yaml                    │
        └─────────────────────────────────────┘
```

## Generic Building Blocks

### 1. Ecosystem Detector
**Purpose**: Identify which languages/build tools are present in a project

**Generic Concepts**:
- Manifest file patterns (e.g., `pom.xml`, `package.json`)
- Directory structure indicators
- File signature detection

**Configuration-Driven**:
- Manifest file names and locations
- Required vs. optional files
- Detection priority/order

### 2. Build Freshness Verifier
**Purpose**: Verify that build artifacts are up-to-date with source manifests

**Generic Concepts**:
- Manifest file (source of truth)
- Cache location (where dependencies are stored)
- Build output location (where artifacts are generated)
- Timestamp comparison logic

**Configuration-Driven**:
- Manifest file path
- Cache directory patterns
- Build output patterns
- Comparison rules (what to check)

### 3. Dependency Auditor
**Purpose**: Verify dependencies are correctly resolved

**Generic Concepts**:
- Dependency declaration (in manifest)
- Dependency resolution (in lock file)
- Dependency artifacts (in cache)

**Configuration-Driven**:
- Manifest parsing rules
- Lock file formats
- Cache structure
- Validation commands

### 4. Environment Variable Auditor
**Purpose**: Check for missing or incorrect environment variables

**Generic Concepts**:
- Variable references (in code)
- Variable declarations (in config files)
- Active environment (system/env)

**Configuration-Driven**:
- Code patterns to search for (regex)
- Config file locations
- Variable naming conventions

### 5. Infrastructure Parity Checker
**Purpose**: Verify required services are running and correct versions

**Generic Concepts**:
- Service name
- Health check endpoint/command
- Version check method
- Expected version/state

**Configuration-Driven**:
- Service detection methods
- Health check commands
- Version extraction patterns

### 6. Environment Reconciler
**Purpose**: Automatically fix detected issues

**Generic Concepts**:
- Issue type
- Fix action
- Verification step

**Configuration-Driven**:
- Issue → fix command mapping
- Pre-fix validation
- Post-fix verification

## Configuration Schema

See `configuration-schema.md` for detailed YAML schema definitions.

## Data Flow

1. **Detection Phase**:
   - Scan project directory for manifest files
   - Load matching ecosystem configs
   - Identify active ecosystems

2. **Verification Phase**:
   - For each detected ecosystem:
     - Load ecosystem-specific config
     - Execute generic checks using config-defined commands
     - Collect results

3. **Reporting Phase**:
   - Aggregate results from all ecosystems
   - Format for MCP response
   - Return to AI agent

4. **Reconciliation Phase** (if requested):
   - Map detected issues to fix commands (from config)
   - Execute fixes
   - Re-verify

## Benefits of This Architecture

1. **Rapid Extension**: Adding Python support = adding `python-pip.yaml` and `python-poetry.yaml`
2. **Community Contributions**: Non-Go developers can contribute ecosystem configs
3. **Version Flexibility**: Different versions of same tool can have different configs
4. **Testing**: Can test core engine with mock configs
5. **Customization**: Users can override/extend configs for their specific setups

