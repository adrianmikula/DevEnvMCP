# Version Management Example

## Java Version Detection Flow

### 1. Detection

**Command**: `java -version 2>&1`

**Output Example**:
```
openjdk version "17.0.9" 2023-10-17
OpenJDK Runtime Environment (build 17.0.9+11)
OpenJDK 64-Bit Server VM (build 17.0.9+11, mixed mode, sharing)
```

**Parsed Result**:
- Language: `java`
- Version: `17.0.9`
- Major: `17`
- Minor: `0`
- Patch: `9`
- Runtime Variant: `OpenJDK` (Provider: `OpenJDK`)

### 2. Validation

**Requirements** (from config):
- Min Version: `11`
- Max Version: `21`
- Preferred Versions: `17`, `21`
- Preferred Runtimes: `Eclipse Temurin`, `OpenJDK`, `Amazon Corretto`

**Validation Result**:
- ✅ Version `17.0.9` is within range (11-21)
- ✅ Version `17` is in preferred list
- ✅ Runtime `OpenJDK` is in preferred list
- **Status**: Valid

### 3. Incompatibility Scenario

**Detected Version**: `8.0.352` (Java 8)

**Validation Result**:
- ❌ Version `8.0.352` is below minimum `11`
- **Issue**: `version_too_old`
- **Suggestion**: Switch to version `17` or `21`

**Suggested Commands** (if sdkman detected):
```bash
sdk install java 17.0.9-tem
sdk use java 17.0.9-tem
```

### 4. Runtime Variant Scenario

**Detected**: `Oracle JDK 17.0.9`

**Requirements**: Preferred runtimes include `Eclipse Temurin`, `OpenJDK`

**Validation Result**:
- ⚠️ Runtime `Oracle JDK` is not in preferred list
- **Issue**: `runtime_not_preferred`
- **Suggestion**: Switch to `Eclipse Temurin` or `OpenJDK`

## Node.js Example

### Configuration

```yaml
version:
  language: "node"
  version_command: "node --version"
  version_pattern: "v(\\d+\\.\\d+\\.\\d+)"
  version_managers:
    - name: "nvm"
      check_command: "command -v nvm"
      list_command: "nvm list"
      install_command: "nvm install {version}"
      switch_command: "nvm use {version}"
      current_command: "nvm current"
    - name: "fnm"
      check_command: "command -v fnm"
      list_command: "fnm list"
      install_command: "fnm install {version}"
      switch_command: "fnm use {version}"
      current_command: "fnm current"

requirements:
  min_version: "18.0.0"
  max_version: "20.0.0"
  preferred_versions:
    - "18.17.0"
    - "20.10.0"
```

### Detection Flow

**Command**: `node --version`

**Output**: `v18.17.0`

**Parsed**:
- Version: `18.17.0`
- Major: `18`
- Valid: ✅ (within 18.0.0-20.0.0)
- Preferred: ✅ (18.17.0 is preferred)

## Integration with Infrastructure Check

When `check_infrastructure_parity` is called:

1. **Detects ecosystem** in project
2. **Checks language version** (if configured)
3. **Validates against requirements**
4. **Reports issues** with suggestions
5. **Checks other services** (Maven, npm, etc.)

**Example Report**:
```
❌ Infrastructure issues found:

- java: Language version incompatibility detected
  - Version 8.0.352 is below minimum required 11
  Suggestion: Switch to a compatible version (required: 11)
    Commands:
      - sdk install java 17.0.9-tem
      - sdk use java 17.0.9-tem
  Available versions: [17, 21]

✅ maven: maven is running (version: 3.9.5)
```

## Benefits

1. **Precise Detection**: Knows exact version and runtime
2. **Actionable Suggestions**: Provides specific commands
3. **Version Manager Aware**: Works with sdkman, nvm, etc.
4. **Runtime Variant Support**: Handles Java runtime variants
5. **Configuration-Driven**: All logic in YAML configs

