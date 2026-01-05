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
├── language-configs/       # Language-level ecosystem configs
├── tool-configs/           # Tool-specific configs organized by language
│   ├── java/              # Java tools (Maven, Gradle, Spring, etc.)
│   ├── python/            # Python tools (Poetry, Conda)
│   ├── javascript/        # JavaScript tools (npm, React, Vite, etc.)
│   ├── csharp/            # C# tools
│   ├── docker/            # Docker tools
│   └── postgres/          # PostgreSQL tools
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

## Adding New Ecosystems

### Adding a Language

Add a YAML config file to `language-configs/` for the base language support.

### Adding a Tool

Add a YAML config file to `tool-configs/{language}/` where `{language}` is the language the tool is used with (e.g., `java`, `python`, `javascript`).

For standalone tools (like Docker, PostgreSQL), add them to `tool-configs/{tool-name}/`.

No code changes needed - the system automatically discovers all configs!

See `docs/architecture/configuration-schema.md` for schema details.

