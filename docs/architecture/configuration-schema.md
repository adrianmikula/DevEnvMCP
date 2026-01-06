# Configuration Schema Reference

This document defines the YAML schema for ecosystem configuration files.

## File Structure

Ecosystem configuration files are organized under a top-level `config/` directory:

### Language Configs (`config/languages/`)

Language-level configurations that define base language support (YAML files directly in this directory):
- `java.yaml` - Java language support
- `python.yaml` - Python language support
- `javascript.yaml` - JavaScript language support
- `csharp.yaml` - C# language support

### Language-Specific Tool Configs (`config/tools/{language}/`)

Tool-specific configurations organized by language in subdirectories:

- `config/tools/java/` - Java tools
  - `maven.yaml` - Maven build tool
  - `gradle.yaml` - Gradle build tool
  - `spring.yaml` - Spring Framework
  - `tomcat.yaml` - Apache Tomcat
  - `jboss.yaml` - JBoss/WildFly

- `config/tools/python/` - Python tools
  - `poetry.yaml` - Poetry package manager
  - `conda.yaml` - Conda environment manager

- `config/tools/javascript/` - JavaScript tools
  - `npm.yaml` - npm package manager
  - `react.yaml` - React framework
  - `vite.yaml` - Vite build tool
  - `webpack.yaml` - Webpack bundler
  - `rollup.yaml` - Rollup bundler
  - `sass.yaml` - Sass/SCSS preprocessor

- `config/tools/csharp/` - C# tools
  - `msbuild.yaml` - MSBuild build tool
  - `nuget.yaml` - NuGet package manager
  - `dotnet.yaml` - dotnet CLI

### Infrastructure Tools (`config/infrastructure/`)

Infrastructure-related tool configurations organized by category:

- `config/infrastructure/docker/` - Docker tools
  - `docker.yaml` - Docker and Docker Compose

- `config/infrastructure/containers/` - Container-related tools
  - (container tool configs can be added here)

- `config/infrastructure/databases/` - Database tools
  - `postgres/postgres.yaml` - PostgreSQL database
  - (other database configs can be added here)

The system automatically discovers all configs recursively from the `config/` directory structure. The loader also supports the old structure (`language-configs/` and `tool-configs/`) for backwards compatibility.

## Root Schema

```yaml
ecosystem:
  name: string              # Human-readable name (e.g., "Java Maven")
  id: string               # Unique identifier (e.g., "java-maven")
  version: string          # Config schema version (e.g., "1.0")
  
  detection:
    # How to detect this ecosystem in a project
    manifest_files: []     # List of manifest file patterns
    required_files: []     # Files that must exist
    optional_files: []    # Files that help confirm presence
    directory_patterns: [] # Directory patterns to look for
    
  manifest:
    # Information about the manifest file
    primary_file: string   # Main manifest file name
    location: string       # Typical location (relative to project root)
    format: string        # Format: "xml", "json", "yaml", "properties", etc.
    
  cache:
    # Dependency cache information
    locations: []          # List of cache directory patterns
    structure: string      # Cache structure type: "flat", "hierarchical", "custom"
    artifact_pattern: string # Pattern to match artifacts (regex)
    
  build:
    # Build output information
    output_directories: [] # Where build artifacts are generated
    artifact_patterns: []  # Patterns for build artifacts
    clean_command: string  # Command to clean build artifacts
    
  dependencies:
    # Dependency management
    lock_file: string      # Lock file name (if applicable)
    lock_file_format: string # Format of lock file
    resolve_command: string # Command to resolve dependencies
    check_command: string  # Command to verify dependency consistency
    
  verification:
    # Commands and patterns for verification
    build_freshness:
      manifest_timestamp_check: boolean
      cache_timestamp_check: boolean
      build_output_check: boolean
      commands: []         # Commands to run for verification
      
    dependency_audit:
      enabled: boolean
      commands: []         # Commands to check dependencies
      
  environment:
    # Environment variable handling
    variable_patterns: []  # Regex patterns to find env var references in code
    config_files: []       # Config files that declare env vars
    required_vars: []      # List of commonly required vars (optional)
    
  infrastructure:
    # Service/infrastructure requirements
    services: []          # List of services this ecosystem typically needs
    
  reconciliation:
    # Auto-fix commands
    fixes:
      - issue_type: string # Type of issue (e.g., "stale_cache", "missing_dep")
        command: string    # Command to fix
        verify_command: string # Command to verify fix worked
        description: string # Human-readable description
```

## Example Configurations

### Java Maven Example

```yaml
ecosystem:
  name: "Java Maven"
  id: "java-maven"
  version: "1.0"
  
  detection:
    manifest_files:
      - "pom.xml"
    required_files:
      - "pom.xml"
    directory_patterns:
      - "src/main/java"
      - "src/test/java"
      
  manifest:
    primary_file: "pom.xml"
    location: "."
    format: "xml"
    
  cache:
    locations:
      - "${HOME}/.m2/repository"
      - "${M2_HOME}/repository"
    structure: "hierarchical"
    artifact_pattern: ".*\\.jar$"
    
  build:
    output_directories:
      - "target/classes"
      - "target/test-classes"
    artifact_patterns:
      - "target/*.jar"
      - "target/*.war"
    clean_command: "mvn clean"
    
  dependencies:
    lock_file: ""  # Maven doesn't use lock files by default
    lock_file_format: ""
    resolve_command: "mvn dependency:resolve"
    check_command: "mvn dependency:tree"
    
  verification:
    build_freshness:
      manifest_timestamp_check: true
      cache_timestamp_check: true
      build_output_check: true
      commands:
        - name: "check_pom_vs_cache"
          type: "timestamp_compare"
          source: "pom.xml"
          target_pattern: "${HOME}/.m2/repository/**/*.jar"
          description: "Compare pom.xml timestamp with cached JARs"
          
        - name: "check_pom_vs_build"
          type: "timestamp_compare"
          source: "pom.xml"
          target_pattern: "target/**/*.class"
          description: "Compare pom.xml timestamp with compiled classes"
          
    dependency_audit:
      enabled: true
      commands:
        - name: "verify_dependencies"
          type: "command"
          command: "mvn dependency:tree"
          parse_output: true
          description: "Verify all dependencies are resolvable"
          
  environment:
    variable_patterns:
      - "\\$\\{([A-Z_][A-Z0-9_]*)\\}"  # ${VAR_NAME}
      - "@Value\\(\"\\$\\{([A-Z_][A-Z0-9_]*)\\}\"\\)"  # @Value("${VAR_NAME}")
      - "System\\.getenv\\(\"([A-Z_][A-Z0-9_]*)\"\\)"  # System.getenv("VAR_NAME")
    config_files:
      - "src/main/resources/application.properties"
      - "src/main/resources/application.yml"
      - "src/main/resources/application-*.yml"
    required_vars: []  # Project-specific
    
  infrastructure:
    services:
      - name: "maven"
        type: "command"
        check_command: "mvn --version"
        version_extract: "Apache Maven (\\d+\\.\\d+\\.\\d+)"
        
  reconciliation:
    fixes:
      - issue_type: "stale_cache"
        command: "mvn clean"
        verify_command: "mvn validate"
        description: "Clean Maven build artifacts"
        
      - issue_type: "missing_dependencies"
        command: "mvn dependency:resolve"
        verify_command: "mvn dependency:tree"
        description: "Resolve missing Maven dependencies"
        
      - issue_type: "corrupted_cache"
        command: "mvn dependency:purge-local-repository -DmanualInclude=*"
        verify_command: "mvn dependency:resolve"
        description: "Purge and re-download corrupted cache entries"
```

### npm Example

```yaml
ecosystem:
  name: "npm"
  id: "npm"
  version: "1.0"
  
  detection:
    manifest_files:
      - "package.json"
    required_files:
      - "package.json"
    optional_files:
      - "package-lock.json"
    directory_patterns:
      - "node_modules"
      
  manifest:
    primary_file: "package.json"
    location: "."
    format: "json"
    
  cache:
    locations:
      - "${HOME}/.npm"
      - "${APPDATA}/npm-cache"  # Windows
    structure: "flat"
    artifact_pattern: ".*"
    
  build:
    output_directories:
      - "dist"
      - "build"
      - "lib"
    artifact_patterns:
      - "dist/**/*"
      - "build/**/*"
    clean_command: "npm run clean || rm -rf dist build"
    
  dependencies:
    lock_file: "package-lock.json"
    lock_file_format: "json"
    resolve_command: "npm install"
    check_command: "npm ls"
    
  verification:
    build_freshness:
      manifest_timestamp_check: true
      cache_timestamp_check: false  # npm cache is less critical
      build_output_check: true
      commands:
        - name: "check_package_vs_lock"
          type: "timestamp_compare"
          source: "package.json"
          target: "package-lock.json"
          description: "Check if package.json is newer than lock file"
          
        - name: "check_lock_vs_node_modules"
          type: "timestamp_compare"
          source: "package-lock.json"
          target: "node_modules"
          description: "Check if node_modules matches lock file"
          
        - name: "check_package_vs_build"
          type: "timestamp_compare"
          source: "package.json"
          target_pattern: "dist/**/*"
          description: "Check if build artifacts are stale"
          
    dependency_audit:
      enabled: true
      commands:
        - name: "verify_dependencies"
          type: "command"
          command: "npm ls --depth=0"
          parse_output: true
          success_pattern: ".*"
          error_pattern: ".*UNMET DEPENDENCY.*|.*extraneous.*"
          description: "Verify all dependencies are installed correctly"
          
  environment:
    variable_patterns:
      - "process\\.env\\.([A-Z_][A-Z0-9_]*)"
      - "\\$\\{([A-Z_][A-Z0-9_]*)\\}"
      - "process\\.env\\[['\"]([A-Z_][A-Z0-9_]*)['\"]\\]"
    config_files:
      - ".env"
      - ".env.local"
      - ".env.*"
    required_vars: []  # Project-specific
    
  infrastructure:
    services:
      - name: "node"
        type: "command"
        check_command: "node --version"
        version_extract: "v(\\d+\\.\\d+\\.\\d+)"
        
      - name: "npm"
        type: "command"
        check_command: "npm --version"
        version_extract: "(\\d+\\.\\d+\\.\\d+)"
        
  reconciliation:
    fixes:
      - issue_type: "stale_lock"
        command: "npm install"
        verify_command: "npm ls --depth=0"
        description: "Update package-lock.json to match package.json"
        
      - issue_type: "missing_dependencies"
        command: "npm install"
        verify_command: "npm ls --depth=0"
        description: "Install missing npm dependencies"
        
      - issue_type: "corrupted_node_modules"
        command: "rm -rf node_modules && npm install"
        verify_command: "npm ls --depth=0"
        description: "Reinstall corrupted node_modules"
        
      - issue_type: "stale_build"
        command: "npm run build"
        verify_command: "test -d dist || test -d build"
        description: "Rebuild stale artifacts"
```

## Command Types

### timestamp_compare
Compare timestamps between source and target files/directories.

```yaml
- name: string
  type: "timestamp_compare"
  source: string           # Source file path
  target: string          # Target file/directory path (or pattern)
  target_pattern: string   # Alternative: glob pattern for targets
  description: string
```

### command
Execute a shell command and parse output.

```yaml
- name: string
  type: "command"
  command: string         # Command to execute
  parse_output: boolean   # Whether to parse command output
  success_pattern: string # Regex for success (optional)
  error_pattern: string   # Regex for errors (optional)
  description: string
```

### file_exists
Check if files exist.

```yaml
- name: string
  type: "file_exists"
  files: []              # List of file paths/patterns
  all_required: boolean  # If true, all must exist; if false, any
  description: string
```

### version_check
Check version of a tool or service.

```yaml
- name: string
  type: "version_check"
  command: string        # Command to get version
  version_extract: string # Regex to extract version from output
  min_version: string    # Minimum required version (optional)
  max_version: string    # Maximum allowed version (optional)
  description: string
```

## Variable Substitution

Configuration files support environment variable substitution:
- `${HOME}` - User home directory
- `${M2_HOME}` - Maven home (if set)
- `${APPDATA}` - Windows AppData directory
- `${PROJECT_ROOT}` - Root of the project being checked

## Configuration Loading Order

1. Load default configs from `ecosystem-configs/` directory
2. Load user overrides from `~/.sentinel/configs/` (if exists)
3. Load project-specific configs from `.sentinel/configs/` (if exists)
4. Project-specific configs override user configs, which override defaults

