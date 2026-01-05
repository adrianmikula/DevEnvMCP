# Example Workflow: How the Architecture Works

This document demonstrates how the language-agnostic architecture works with a concrete example.

## Scenario: Checking a Java Maven Project

### Step 1: Project Detection

**User Action**: AI agent calls `verify_build_freshness` tool with project path `/my-project`

**System Action**:
1. Core engine scans `/my-project` directory
2. Finds `pom.xml` file
3. Matches against ecosystem configs:
   - `java-maven.yaml` matches (has `pom.xml` in `required_files`)
   - `java-gradle.yaml` doesn't match (no `build.gradle`)
   - `npm.yaml` doesn't match (no `package.json`)
4. Returns: `DetectedEcosystem{ID: "java-maven", Confidence: 1.0}`

**Code Flow**:
```go
// detector/detector.go
ecosystems := DetectEcosystems("/my-project", allConfigs)
// Returns: [DetectedEcosystem{ID: "java-maven", Config: <loaded config>}]
```

### Step 2: Configuration Loading

**System Action**:
1. Load `ecosystem-configs/java-maven.yaml`
2. Parse into `EcosystemConfig` struct
3. Apply any user/project overrides (if present)
4. Validate configuration

**Configuration Used** (excerpt):
```yaml
ecosystem:
  id: "java-maven"
  verification:
    build_freshness:
      commands:
        - name: "check_pom_vs_cache"
          type: "timestamp_compare"
          source: "pom.xml"
          target_pattern: "${HOME}/.m2/repository/**/*.jar"
```

**Code Flow**:
```go
// config/loader.go
config := LoadEcosystemConfig("java-maven")
// Returns: *EcosystemConfig with all verification commands loaded
```

### Step 3: Build Freshness Verification

**System Action**:
For each verification command in config:

**Command 1: `check_pom_vs_cache`**
1. Get timestamp of `pom.xml`: `2026-01-15 10:30:00`
2. Expand pattern: `${HOME}/.m2/repository/**/*.jar`
3. Find matching JARs in cache
4. Get newest JAR timestamp: `2026-01-15 09:00:00`
5. Compare: `pom.xml` is newer → **STALE CACHE detected**

**Command 2: `check_pom_vs_build`**
1. Get timestamp of `pom.xml`: `2026-01-15 10:30:00`
2. Find compiled classes: `target/classes/**/*.class`
3. Get newest class timestamp: `2026-01-15 08:00:00`
4. Compare: `pom.xml` is newer → **STALE BUILD detected**

**Code Flow**:
```go
// verifier/freshness.go
report := VerifyBuildFreshness("/my-project", ecosystem)

// For each command in config.Verification.BuildFreshness.Commands:
for _, cmd := range config.Verification.BuildFreshness.Commands {
    switch cmd.Type {
    case "timestamp_compare":
        result := CompareTimestamps(
            resolvePath(cmd.Source, projectRoot),
            resolvePattern(cmd.TargetPattern, projectRoot),
        )
        // result: {Stale: true, SourceTime: ..., TargetTime: ...}
    }
}
```

### Step 4: Result Aggregation

**System Action**:
1. Collect all verification results
2. Aggregate into structured report
3. Format for MCP response

**Result**:
```json
{
  "ecosystem": "java-maven",
  "status": "issues_found",
  "issues": [
    {
      "type": "stale_cache",
      "severity": "warning",
      "message": "pom.xml is newer than cached JARs in ~/.m2/repository",
      "fix_available": true,
      "fix_command": "mvn clean"
    },
    {
      "type": "stale_build",
      "severity": "error",
      "message": "pom.xml is newer than compiled classes in target/",
      "fix_available": true,
      "fix_command": "mvn clean compile"
    }
  ]
}
```

### Step 5: MCP Response

**System Action**:
Format result as MCP tool response and return to AI agent

**MCP Response**:
```json
{
  "content": [
    {
      "type": "text",
      "text": "Build Freshness Report for java-maven:\n\n⚠️  Stale Cache: pom.xml is newer than cached JARs\n   Fix: mvn clean\n\n❌ Stale Build: pom.xml is newer than compiled classes\n   Fix: mvn clean compile"
    }
  ]
}
```

## Scenario: Adding Python Support (No Code Changes!)

### Step 1: Create Configuration File

**Action**: Create `ecosystem-configs/python-pip.yaml`

```yaml
ecosystem:
  name: "Python pip"
  id: "python-pip"
  version: "1.0"
  
  detection:
    manifest_files:
      - "requirements.txt"
      - "setup.py"
      - "pyproject.toml"
    required_files:
      - "requirements.txt"
      
  manifest:
    primary_file: "requirements.txt"
    location: "."
    format: "text"
    
  cache:
    locations:
      - "${HOME}/.cache/pip"
    structure: "flat"
    
  build:
    output_directories:
      - "dist"
      - "build"
      - "*.egg-info"
    clean_command: "pip uninstall -y -r requirements.txt || true"
    
  verification:
    build_freshness:
      commands:
        - name: "check_requirements_vs_cache"
          type: "timestamp_compare"
          source: "requirements.txt"
          target_pattern: "${HOME}/.cache/pip/**/*"
          description: "Check if requirements.txt is newer than pip cache"
          
        - name: "check_requirements_vs_build"
          type: "timestamp_compare"
          source: "requirements.txt"
          target_pattern: "dist/**/*"
          description: "Check if build artifacts are stale"
```

### Step 2: Test

**Action**: Run Sentinel on a Python project

**Result**: 
- Python project detected automatically
- Build freshness checks run using pip-specific commands
- **No code changes needed!**

## Scenario: Custom Project Configuration

### Step 1: Project-Specific Override

**Action**: Create `.sentinel/configs/java-maven.yaml` in project root

```yaml
ecosystem:
  id: "java-maven"
  
  # Override cache location for this project
  cache:
    locations:
      - "/custom/maven/repo"  # Custom Maven repo location
      
  # Add custom verification
  verification:
    build_freshness:
      commands:
        - name: "check_custom_artifacts"
          type: "timestamp_compare"
          source: "pom.xml"
          target_pattern: "/custom/artifacts/**/*.jar"
          description: "Check custom artifact location"
```

### Step 2: System Behavior

**System Action**:
1. Load default `java-maven.yaml`
2. Load project-specific `.sentinel/configs/java-maven.yaml`
3. Merge configs (project overrides default)
4. Use merged config for verification

**Result**: Project uses custom Maven repo location and additional checks

## Key Benefits Demonstrated

1. **Language-Agnostic Core**: Same Go code works for Java, npm, Python, etc.
2. **Easy Extension**: Adding Python = adding one YAML file
3. **Flexibility**: Projects can override/extend configs
4. **Maintainability**: Tool-specific changes don't require code changes
5. **Testability**: Can test core engine with mock configs

## Command Execution Flow

### Example: Executing a Fix Command

**User Action**: AI agent calls `reconcile_environment` with issue type `stale_cache`

**System Action**:
1. Look up fix in config:
   ```yaml
   reconciliation:
     fixes:
       - issue_type: "stale_cache"
         command: "mvn clean"
         verify_command: "mvn validate"
   ```
2. Execute: `mvn clean` in project directory
3. Verify: Execute `mvn validate` to confirm fix
4. Return result to AI agent

**Code Flow**:
```go
// reconciler/reconciler.go
fix := config.Reconciliation.Fixes.FindByIssueType("stale_cache")
result := ExecuteFix(fix, projectRoot)

// reconciler/fix_executor.go
func ExecuteFix(fix *FixConfig, projectRoot string) (*FixResult, error) {
    // Execute fix command
    cmd := exec.Command("sh", "-c", fix.Command)
    cmd.Dir = projectRoot
    err := cmd.Run()
    
    if err != nil {
        return &FixResult{Success: false, Error: err}, nil
    }
    
    // Verify fix
    verifyCmd := exec.Command("sh", "-c", fix.VerifyCommand)
    verifyCmd.Dir = projectRoot
    verifyErr := verifyCmd.Run()
    
    return &FixResult{
        Success: verifyErr == nil,
        Message: fix.Description,
    }, nil
}
```

## Summary

This architecture enables:
- **Zero code changes** for new ecosystems (just add YAML)
- **Consistent behavior** across all ecosystems (same core logic)
- **Project customization** (override configs per project)
- **Easy testing** (mock configs for unit tests)
- **Community contributions** (non-Go devs can add ecosystem configs)

