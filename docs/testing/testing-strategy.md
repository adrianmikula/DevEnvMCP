# Testing Strategy for Dev-Env Sentinel MCP

This document outlines the testing approach, tools, and best practices for the Dev-Env Sentinel MCP project.

## Testing Philosophy

1. **Test-Driven Development (TDD)**: Write tests before implementation when possible
2. **Comprehensive Coverage**: Aim for high coverage of critical paths (80%+)
3. **Fast Execution**: Tests should run quickly (< 1 second for unit tests)
4. **Isolation**: Each test should be independent and runnable in parallel
5. **Readability**: Tests should clearly express what they're testing

## Testing Stack

### Core Testing Tools

#### 1. **Go Standard Library `testing` Package**
- **Purpose**: Foundation for all tests
- **Usage**: Built-in, no installation needed
- **Features**: Unit tests, benchmarks, subtests, parallel execution

#### 2. **testify** (Recommended)
- **Package**: `github.com/stretchr/testify`
- **Purpose**: Assertions, mocking, and test suites
- **Install**: `go get github.com/stretchr/testify`
- **Components**:
  - `assert`: Rich assertion library
  - `require`: Same as assert but stops test on failure
  - `mock`: Mock object generation
  - `suite`: Test suite support

#### 3. **gomock** (For Interface Mocking)
- **Package**: `go.uber.org/mock` (or `github.com/golang/mock`)
- **Purpose**: Generate mocks from interfaces
- **Install**: `go install go.uber.org/mock/mockgen@latest`
- **Usage**: Code generation for mock interfaces

### Optional Testing Tools

#### 4. **httptest** (Standard Library)
- **Purpose**: Testing HTTP handlers and clients
- **Usage**: For MCP server HTTP transport testing

#### 5. **testcontainers-go** (For Integration Tests)
- **Package**: `github.com/testcontainers/testcontainers-go`
- **Purpose**: Spin up Docker containers for integration tests
- **Use Case**: Testing with real Maven/npm projects in containers

## Testing Patterns

### 1. Table-Driven Tests (Go Idiom)

**When to Use**: Testing multiple input/output scenarios

**Example**:
```go
func TestDetectEcosystem(t *testing.T) {
    tests := []struct {
        name        string
        projectRoot string
        files       []string
        want        []string
        wantErr     bool
    }{
        {
            name:        "detects maven project",
            projectRoot: "/tmp/test",
            files:       []string{"pom.xml"},
            want:        []string{"java-maven"},
            wantErr:     false,
        },
        {
            name:        "detects npm project",
            projectRoot: "/tmp/test",
            files:       []string{"package.json"},
            want:        []string{"npm"},
            wantErr:     false,
        },
        {
            name:        "detects multiple ecosystems",
            projectRoot: "/tmp/test",
            files:       []string{"pom.xml", "package.json"},
            want:        []string{"java-maven", "npm"},
            wantErr:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup: Create test directory with files
            tmpDir := setupTestProject(t, tt.projectRoot, tt.files)
            defer os.RemoveAll(tmpDir)

            // Execute
            got, err := DetectEcosystems(tmpDir, allConfigs)

            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("DetectEcosystems() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            assert.Equal(t, tt.want, ecosystemIDs(got))
        })
    }
}
```

### 2. Test Fixtures and Temporary Files

**When to Use**: Testing file operations, config loading, filesystem scanning

**Example**:
```go
func TestLoadEcosystemConfig(t *testing.T) {
    // Create temporary config file
    tmpDir := t.TempDir() // Go 1.15+ - automatically cleaned up
    configPath := filepath.Join(tmpDir, "java-maven.yaml")
    
    // Write test config
    testConfig := `
ecosystem:
  id: "java-maven"
  detection:
    required_files:
      - "pom.xml"
`
    err := os.WriteFile(configPath, []byte(testConfig), 0644)
    require.NoError(t, err)

    // Test
    config, err := LoadEcosystemConfig(configPath)
    require.NoError(t, err)
    assert.Equal(t, "java-maven", config.ID)
}
```

**Using `testdata/` Directory**:
```go
// testdata/java-maven-project/
//   ├── pom.xml
//   └── src/main/java/...

func TestVerifyBuildFreshness_RealProject(t *testing.T) {
    projectRoot := filepath.Join("testdata", "java-maven-project")
    
    ecosystems, err := DetectEcosystems(projectRoot, allConfigs)
    require.NoError(t, err)
    require.Len(t, ecosystems, 1)
    
    report, err := VerifyBuildFreshness(projectRoot, ecosystems[0])
    require.NoError(t, err)
    assert.NotNil(t, report)
}
```

### 3. Mocking Dependencies

**When to Use**: Isolating units under test, avoiding external dependencies

#### Using Interfaces for Testability

```go
// Define interface for command execution
type CommandExecutor interface {
    Execute(command string, dir string) (*CommandResult, error)
}

// Real implementation
type RealCommandExecutor struct{}

func (e *RealCommandExecutor) Execute(cmd string, dir string) (*CommandResult, error) {
    // Real os/exec implementation
}

// Test implementation
type MockCommandExecutor struct {
    mock.Mock
}

func (m *MockCommandExecutor) Execute(cmd string, dir string) (*CommandResult, error) {
    args := m.Called(cmd, dir)
    return args.Get(0).(*CommandResult), args.Error(1)
}

// Usage in tests
func TestVerifyBuildFreshness_MockedExecutor(t *testing.T) {
    mockExecutor := new(MockCommandExecutor)
    
    mockExecutor.On("Execute", "mvn validate", "/project").Return(
        &CommandResult{Success: true, Output: "BUILD SUCCESS"},
        nil,
    )
    
    verifier := NewVerifier(mockExecutor)
    result, err := verifier.VerifyBuildFreshness("/project", ecosystem)
    
    require.NoError(t, err)
    assert.True(t, result.IsHealthy)
    mockExecutor.AssertExpectations(t)
}
```

#### Using gomock for Generated Mocks

```go
//go:generate mockgen -source=executor.go -destination=mock_executor_test.go

// In test file
func TestWithGeneratedMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockExecutor := NewMockCommandExecutor(ctrl)
    mockExecutor.EXPECT().
        Execute("mvn validate", "/project").
        Return(&CommandResult{Success: true}, nil).
        Times(1)
    
    // Use mockExecutor in test
}
```

### 4. Testing Command Execution

**Challenge**: Testing `os/exec` calls without actually running commands

**Solution**: Use interfaces and dependency injection

```go
// executor.go
type Executor interface {
    Run(ctx context.Context, cmd *exec.Cmd) error
    Output(ctx context.Context, cmd *exec.Cmd) ([]byte, error)
}

type RealExecutor struct{}

func (e *RealExecutor) Run(ctx context.Context, cmd *exec.Cmd) error {
    return cmd.Run()
}

// executor_test.go
type FakeExecutor struct {
    commands map[string]*CommandResult
}

func (f *FakeExecutor) Run(ctx context.Context, cmd *exec.Cmd) error {
    key := strings.Join(cmd.Args, " ")
    if result, ok := f.commands[key]; ok {
        if !result.Success {
            return result.Error
        }
        return nil
    }
    return fmt.Errorf("unexpected command: %s", key)
}

func TestExecuteFix(t *testing.T) {
    fakeExecutor := &FakeExecutor{
        commands: map[string]*CommandResult{
            "mvn clean": {Success: true},
        },
    }
    
    reconciler := NewReconciler(fakeExecutor)
    result, err := reconciler.ExecuteFix(fixConfig, "/project")
    
    require.NoError(t, err)
    assert.True(t, result.Success)
}
```

### 5. Testing File Operations

**Pattern**: Use temporary directories and `testdata/`

```go
func TestTimestampCompare(t *testing.T) {
    tmpDir := t.TempDir()
    
    // Create source file
    sourceFile := filepath.Join(tmpDir, "source.txt")
    err := os.WriteFile(sourceFile, []byte("source"), 0644)
    require.NoError(t, err)
    
    // Wait a bit to ensure different timestamps
    time.Sleep(10 * time.Millisecond)
    
    // Create target file (newer)
    targetFile := filepath.Join(tmpDir, "target.txt")
    err = os.WriteFile(targetFile, []byte("target"), 0644)
    require.NoError(t, err)
    
    // Test
    isNewer, err := CompareTimestamps(sourceFile, targetFile)
    require.NoError(t, err)
    assert.False(t, isNewer) // target is newer
}
```

### 6. Testing Configuration Loading

```go
func TestLoadEcosystemConfig(t *testing.T) {
    tests := []struct {
        name      string
        configYAML string
        wantErr   bool
        validate  func(t *testing.T, config *EcosystemConfig)
    }{
        {
            name: "valid config",
            configYAML: `
ecosystem:
  id: "java-maven"
  name: "Java Maven"
`,
            wantErr: false,
            validate: func(t *testing.T, config *EcosystemConfig) {
                assert.Equal(t, "java-maven", config.ID)
                assert.Equal(t, "Java Maven", config.Name)
            },
        },
        {
            name:      "invalid yaml",
            configYAML: `invalid: yaml: [`,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tmpFile := filepath.Join(t.TempDir(), "config.yaml")
            err := os.WriteFile(tmpFile, []byte(tt.configYAML), 0644)
            require.NoError(t, err)

            config, err := LoadEcosystemConfig(tmpFile)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            if tt.validate != nil {
                tt.validate(t, config)
            }
        })
    }
}
```

## Test Organization

### Directory Structure

```
internal/
├── config/
│   ├── loader.go
│   ├── loader_test.go
│   ├── schema.go
│   └── schema_test.go
├── detector/
│   ├── detector.go
│   └── detector_test.go
└── verifier/
    ├── freshness.go
    └── freshness_test.go

testdata/
├── java-maven-project/
│   ├── pom.xml
│   └── src/
├── npm-project/
│   ├── package.json
│   └── src/
└── multi-ecosystem-project/
    ├── pom.xml
    └── package.json

integration/
└── integration_test.go
```

### Test File Naming

- Unit tests: `*_test.go` in same package
- Integration tests: `*_integration_test.go` or in `integration/` directory
- Benchmarks: `*_bench_test.go` or use `Benchmark*` prefix

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestDetectEcosystem ./internal/detector

# Run tests in parallel
go test -parallel 4 ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Flags

```bash
# Short timeout for fast feedback
go test -timeout 30s ./...

# Count: run each test N times
go test -count=3 ./...

# Race detector (for concurrent code)
go test -race ./...
```

## Integration Testing

### Strategy

1. **Unit Tests**: Fast, isolated, mocked dependencies
2. **Integration Tests**: Slower, test with real filesystems, optional real commands
3. **End-to-End Tests**: Full MCP server with real projects

### Integration Test Example

```go
// integration/integration_test.go
// +build integration

package integration

import (
    "testing"
    "path/filepath"
)

func TestEndToEnd_JavaMavenProject(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    projectRoot := filepath.Join("testdata", "java-maven-project")
    
    // Detect ecosystems
    ecosystems, err := DetectEcosystems(projectRoot, allConfigs)
    require.NoError(t, err)
    require.Len(t, ecosystems, 1)
    
    // Verify build freshness
    report, err := VerifyBuildFreshness(projectRoot, ecosystems[0])
    require.NoError(t, err)
    
    // Assertions
    assert.NotNil(t, report)
    // ... more assertions
}
```

Run with: `go test -tags=integration ./integration`

## Benchmarking

### Example

```go
func BenchmarkDetectEcosystems(b *testing.B) {
    projectRoot := filepath.Join("testdata", "large-project")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = DetectEcosystems(projectRoot, allConfigs)
    }
}
```

Run with: `go test -bench=. -benchmem ./...`

## Best Practices for This Project

### 1. Test Configuration Loading Extensively

- Valid configs
- Invalid YAML
- Missing required fields
- Config merging (default → user → project)

### 2. Test Ecosystem Detection

- Single ecosystem projects
- Multi-ecosystem projects
- Projects with no ecosystems
- Edge cases (empty directories, symlinks)

### 3. Test File Operations Safely

- Always use `t.TempDir()` for temporary files
- Use `testdata/` for fixtures
- Clean up after tests (automatic with `t.TempDir()`)

### 4. Mock External Dependencies

- Command execution (`os/exec`)
- File system operations (if complex)
- Network calls (if any)

### 5. Test Error Handling

- Invalid inputs
- Missing files
- Permission errors
- Timeout scenarios

### 6. Test Configuration-Driven Logic

- Different ecosystem configs
- Custom project configs
- Config overrides

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## Testing Checklist

- [ ] Unit tests for all core functions
- [ ] Table-driven tests for multiple scenarios
- [ ] Integration tests for end-to-end flows
- [ ] Mock external dependencies
- [ ] Test error handling
- [ ] Test edge cases
- [ ] Benchmark critical paths
- [ ] Achieve 80%+ code coverage
- [ ] Tests run in parallel
- [ ] Tests are independent
- [ ] CI/CD integration

## Recommended Libraries Summary

| Library | Purpose | When to Use |
|---------|---------|-------------|
| `testing` (stdlib) | Core testing | Always |
| `testify/assert` | Assertions | Recommended for all tests |
| `testify/require` | Assertions (fatal) | When test can't continue |
| `testify/mock` | Manual mocks | Simple mocking needs |
| `gomock` | Generated mocks | Complex interface mocking |
| `httptest` (stdlib) | HTTP testing | MCP HTTP transport |
| `testcontainers-go` | Docker containers | Integration tests with real tools |

## Next Steps

1. Set up `go.mod` with testing dependencies
2. Create initial test structure
3. Write tests for configuration loading
4. Write tests for ecosystem detection
5. Set up CI/CD pipeline

