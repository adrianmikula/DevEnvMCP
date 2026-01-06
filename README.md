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
├── config/                 # Configuration files
│   ├── languages/         # Language-level configs
│   │   └── *.yaml        # Language configs (java.yaml, python.yaml, etc.)
│   ├── tools/             # Language-specific tool configs
│   │   ├── java/         # Java tools (Maven, Gradle, Spring, etc.)
│   │   ├── python/       # Python tools (Poetry, Conda)
│   │   ├── javascript/   # JavaScript tools (npm, React, Vite, etc.)
│   │   └── csharp/       # C# tools (MSBuild, NuGet, dotnet CLI)
│   └── infrastructure/    # Infrastructure tools
│       ├── docker/        # Docker tools
│       ├── containers/    # Container-related tools
│       └── databases/    # Database tools
│           └── postgres/ # PostgreSQL tools
└── docs/                   # Documentation
```

## Core Principles

1. **DRY (Don't Repeat Yourself)**: Shared utilities in `internal/common`
2. **KISS (Keep It Simple)**: Focused functions, clear responsibilities
3. **Efficiency**: Minimal allocations, optimized file operations
4. **Extensibility**: Add new ecosystems via YAML configs

## Installation

### Via npm/npx (Recommended)

Install globally:
```bash
npm install -g dev-env-sentinel
```

Or use with npx (no installation required):
```bash
npx dev-env-sentinel
```

### For MCP Clients

#### Cursor

Add to your Cursor MCP settings (`~/.cursor/mcp.json` on Linux/Mac, `%APPDATA%\Cursor\mcp.json` on Windows):

```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "npx",
      "args": ["-y", "dev-env-sentinel"]
    }
  }
}
```

Or if installed globally:
```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "dev-env-sentinel"
    }
  }
}
```

#### Claude Code

Add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "npx",
      "args": ["-y", "dev-env-sentinel"]
    }
  }
}
```

#### Google Antigravity

Add to your Antigravity MCP configuration:

```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "npx",
      "args": ["-y", "dev-env-sentinel"]
    }
  }
}
```

## Development

### Prerequisites

- Go 1.13+ (Go Modules support)
- Node.js 14+ (for build scripts)
- YAML configs in `config/` directory structure

### Build

Build for current platform:
```bash
npm run build
```

Build for all platforms (Windows and Linux):
```bash
npm run build:all
```

Build for specific platform:
```bash
npm run build:windows
npm run build:linux
```

Or build directly with Go:
```bash
go build ./cmd/sentinel
```

### Run

Run the MCP server (no arguments for MCP mode):
```bash
./sentinel
```

Or via npm:
```bash
npx dev-env-sentinel
```

## Supported Ecosystems

The following ecosystems are currently supported:

### Phase 1 (MVP)
- **Java**: Maven (`pom.xml`), Gradle (`build.gradle`)
- **npm**: npm, Yarn, pnpm (`package.json`)

### Phase 2 (Implemented)
- **React**: React applications (`package.json` with React)
- **Vite**: Vite build tool (`vite.config.js/ts`)
- **Python**: pip (`requirements.txt`, `setup.py`)
- **Poetry**: Poetry package manager (`pyproject.toml`, `poetry.lock`)
- **Conda**: Conda environment manager (`environment.yml`)
- **Docker**: Docker and Docker Compose (`Dockerfile`, `docker-compose.yml`)
- **PostgreSQL**: Database configuration and migrations
- **C# (.NET)**: .NET projects (`*.csproj`, `*.sln`)

### Build Tools & Frameworks
- **Webpack**: Webpack bundler (`webpack.config.js`)
- **Rollup**: Rollup bundler (`rollup.config.js`)
- **Sass/SCSS**: Sass preprocessor (`*.scss`, `*.sass`)
- **Spring Framework**: Spring Boot/Spring Framework (`application.properties`, `pom.xml`)
- **Apache Tomcat**: Tomcat servlet container (`web.xml`, `context.xml`)
- **JBoss/WildFly**: JBoss application server (`jboss-web.xml`, `standalone.xml`)

## Monetization & Licensing

Dev-Env Sentinel uses a freemium model with feature flags:

- **Free Tier**: Basic verification and auditing tools
- **Pro Tier**: Auto-fix capabilities and advanced features
- **Enterprise Tier**: Docker orchestration and custom configurations

### Quick Start

Check your license status:
```bash
# Via MCP tool
check_license_status()
```

Get Pro license information:
```bash
# Via MCP tool
get_pro_license()
```

Activate a license:
```bash
# Via MCP tool
activate_pro(license_key="your-license-key")
```

### Payment Options

1. **Stripe Payment Link** - One-time or subscription payments
2. **Apify Pay-Per-Event** - $0.02-$0.05 per tool call

See `docs/monetization.md` for detailed information.

## Adding New Ecosystems

### Adding a Language

Add a YAML config file to `config/languages/` for the base language support (e.g., `config/languages/go.yaml`).

### Adding a Tool

For language-specific tools, add a YAML config file to `config/tools/{language}/` where `{language}` is the language the tool is used with (e.g., `config/tools/java/maven.yaml`).

For infrastructure tools, add them to the appropriate subdirectory under `config/infrastructure/`:
- Docker tools: `config/infrastructure/docker/`
- Container tools: `config/infrastructure/containers/`
- Database tools: `config/infrastructure/databases/{tool-name}/`

No code changes needed - the system automatically discovers all configs!

See `docs/architecture/configuration-schema.md` for schema details.

