# Dev-Env Sentinel MCP

A lean, efficient MCP server for monitoring and verifying development environment health.

## Architecture

- **Language-Agnostic**: Core engine works with any language/tool via YAML configs
- **Configuration-Driven**: All tool-specific logic in YAML files
- **DRY/KISS**: Shared utilities, minimal code duplication
- **Efficient**: Optimized for performance and low memory usage

## Project Structure

```
dev-env-sentinel/
├── cmd/
│   └── sentinel/          # MCP server entry point
├── internal/
│   ├── common/            # Shared utilities (DRY)
│   ├── config/            # Configuration loading
│   ├── detector/           # Ecosystem detection
│   ├── verifier/           # Build freshness verification
│   ├── auditor/            # Dependency/env var auditing
│   └── reconciler/         # Auto-fix functionality
├── ecosystem-configs/      # YAML configs for each ecosystem
└── docs/                   # Documentation
```

## Core Principles

1. **DRY (Don't Repeat Yourself)**: Shared utilities in `internal/common`
2. **KISS (Keep It Simple)**: Focused functions, clear responsibilities
3. **Efficiency**: Minimal allocations, optimized file operations
4. **Extensibility**: Add new ecosystems via YAML configs

## Development

### Prerequisites

- Go 1.13+ (Go Modules support)
- YAML configs in `ecosystem-configs/`

### Build

```bash
go build ./cmd/sentinel
```

### Run

```bash
./sentinel <project-root>
```

## Adding New Ecosystems

Simply add a YAML config file to `ecosystem-configs/` - no code changes needed!

See `docs/architecture/configuration-schema.md` for schema details.

