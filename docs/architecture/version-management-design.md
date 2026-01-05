# Version Management Design

## Overview

Version management is critical for detecting and resolving runtime compatibility issues. This system needs to handle:
- Core language version detection
- Version management tool awareness
- Variant runtime detection (especially Java)
- Incompatibility detection
- Actionable suggestions

## Design Principles

1. **Configuration-Driven**: Version info defined in ecosystem configs
2. **Tool-Agnostic**: Support multiple version managers per language
3. **Runtime-Aware**: Detect and distinguish runtime variants
4. **Actionable**: Provide specific commands to fix issues

## Core Concepts

### Version Categories
- **Major Version**: Java 17, Node.js 18, Python 3.11
- **Minor Version**: Java 17.0.9, Node.js 18.17.0
- **Runtime Variant**: OpenJDK, Oracle JDK, Eclipse Temurin, Amazon Corretto

### Version Management Tools
- **Java**: sdkman, jenv, asdf
- **Node.js**: nvm, fnm, asdf
- **Python**: pyenv, asdf, conda
- **System**: Direct installation (no manager)

### Runtime Variants (Java Example)
- **OpenJDK**: Open source reference implementation
- **Oracle JDK**: Oracle's commercial implementation
- **Eclipse Temurin**: Adoptium's OpenJDK builds
- **Amazon Corretto**: Amazon's OpenJDK distribution
- **Azul Zulu**: Azul's OpenJDK builds
- **Microsoft Build**: Microsoft's OpenJDK builds

## Architecture

### Components

1. **Version Detector** (`internal/version/detector.go`)
   - Detects current language version
   - Detects runtime variant (if applicable)
   - Detects version manager in use

2. **Version Manager** (`internal/version/manager.go`)
   - Lists available versions
   - Checks if version is installed
   - Provides install/switch commands

3. **Version Validator** (`internal/version/validator.go`)
   - Validates version against requirements
   - Detects incompatibilities
   - Suggests alternatives

4. **Runtime Detector** (`internal/version/runtime.go`)
   - Detects runtime variant (Java-specific)
   - Maps variant to provider
   - Checks variant compatibility

## Configuration Schema Extension

### New Config Sections

```yaml
version:
  language: "java"  # Language identifier
  version_manager_tools: []  # Available version managers
  runtime_variants: []  # Runtime variants (Java-specific)
  version_pattern: ""  # Regex to extract version
  runtime_pattern: ""  # Regex to extract runtime info
  
requirements:
  min_version: ""  # Minimum required version
  max_version: ""  # Maximum allowed version
  preferred_runtimes: []  # Preferred runtime variants
```

## Detection Strategy

### 1. Language Version Detection
- Execute version command (e.g., `java -version`, `node --version`)
- Parse output using regex pattern from config
- Extract major.minor.patch version

### 2. Runtime Variant Detection (Java)
- Parse `java -version` output for vendor info
- Match against known runtime patterns
- Identify provider (OpenJDK, Oracle, Temurin, etc.)

### 3. Version Manager Detection
- Check for version manager presence
- Detect which manager is active
- List available versions via manager

### 4. Incompatibility Detection
- Compare detected version with requirements
- Check runtime variant compatibility
- Identify specific incompatibility reason

## Suggestion Generation

When incompatibility detected:
1. **Identify Issue**: Version too old/new, wrong runtime, etc.
2. **List Alternatives**: Available versions that meet requirements
3. **Provide Commands**: Specific commands to install/switch
4. **Runtime Options**: Alternative runtime variants if applicable

## Example: Java Version Management

### Detection Flow
1. Run `java -version`
2. Parse: "openjdk version \"17.0.9\" 2023-10-17"
3. Extract: version=17.0.9, runtime=OpenJDK
4. Check requirements: min=17, max=21
5. If incompatible: suggest Java 17-21, provide sdkman/jenv commands

### Runtime Variant Handling
- Detect: "Eclipse Temurin(TM) 17.0.9+11"
- Identify: Eclipse Temurin (Adoptium)
- Check compatibility: Temurin is compatible
- If incompatible: suggest alternative runtimes

## Implementation Plan

1. **Version Detection** (Phase 1)
   - Basic version detection
   - Pattern-based parsing
   - Runtime variant detection (Java)

2. **Version Manager Integration** (Phase 2)
   - Detect version managers
   - List available versions
   - Generate switch commands

3. **Validation & Suggestions** (Phase 3)
   - Compare with requirements
   - Generate suggestions
   - Provide fix commands

4. **Integration** (Phase 4)
   - Add to infrastructure checker
   - Add to MCP tools
   - Update ecosystem configs

