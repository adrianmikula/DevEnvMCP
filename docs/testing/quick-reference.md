# Testing Quick Reference

Quick reference guide for common testing patterns in the Dev-Env Sentinel MCP project.

## Installation

```bash
# Install testify
go get github.com/stretchr/testify

# Install gomock (optional)
go install go.uber.org/mock/mockgen@latest
```

## Common Patterns

### 1. Basic Test with testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFunctionName(t *testing.T) {
    // Setup
    input := "test"
    
    // Execute
    result, err := FunctionName(input)
    
    // Assert
    require.NoError(t, err)  // Stops test on failure
    assert.Equal(t, "expected", result)
}
```

### 2. Table-Driven Test

```go
func TestDetectEcosystem(t *testing.T) {
    tests := []struct {
        name    string
        files   []string
        want    string
        wantErr bool
    }{
        {"maven", []string{"pom.xml"}, "java-maven", false},
        {"npm", []string{"package.json"}, "npm", false},
        {"none", []string{}, "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tmpDir := setupTestDir(t, tt.files)
            defer os.RemoveAll(tmpDir)
            
            got, err := DetectEcosystem(tmpDir)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### 3. Temporary Files/Directories

```go
func TestWithTempDir(t *testing.T) {
    // Go 1.15+ - automatically cleaned up
    tmpDir := t.TempDir()
    
    file := filepath.Join(tmpDir, "test.txt")
    err := os.WriteFile(file, []byte("content"), 0644)
    require.NoError(t, err)
    
    // Test file operations
    // tmpDir automatically deleted after test
}
```

### 4. Mock with testify/mock

```go
// Define interface
type Executor interface {
    Execute(cmd string) error
}

// Create mock
type MockExecutor struct {
    mock.Mock
}

func (m *MockExecutor) Execute(cmd string) error {
    args := m.Called(cmd)
    return args.Error(0)
}

// Use in test
func TestWithMock(t *testing.T) {
    mockExec := new(MockExecutor)
    mockExec.On("Execute", "mvn clean").Return(nil)
    
    result := UseExecutor(mockExec)
    
    assert.NotNil(t, result)
    mockExec.AssertExpectations(t)
}
```

### 5. Testing File Operations

```go
func TestFileOperations(t *testing.T) {
    tmpDir := t.TempDir()
    
    // Create test file
    testFile := filepath.Join(tmpDir, "test.txt")
    err := os.WriteFile(testFile, []byte("content"), 0644)
    require.NoError(t, err)
    
    // Test reading
    content, err := os.ReadFile(testFile)
    require.NoError(t, err)
    assert.Equal(t, "content", string(content))
    
    // Test file exists
    _, err = os.Stat(testFile)
    assert.NoError(t, err)
}
```

### 6. Testing YAML Config Loading

```go
func TestLoadConfig(t *testing.T) {
    configYAML := `
ecosystem:
  id: "test"
  name: "Test Ecosystem"
`
    
    tmpFile := filepath.Join(t.TempDir(), "config.yaml")
    err := os.WriteFile(tmpFile, []byte(configYAML), 0644)
    require.NoError(t, err)
    
    config, err := LoadEcosystemConfig(tmpFile)
    require.NoError(t, err)
    assert.Equal(t, "test", config.ID)
}
```

### 7. Testing Command Execution (Mocked)

```go
// Define executor interface
type CmdRunner interface {
    Run(ctx context.Context, name string, args ...string) error
}

// Fake implementation for tests
type FakeCmdRunner struct {
    commands map[string]error
}

func (f *FakeCmdRunner) Run(ctx context.Context, name string, args ...string) error {
    key := name + " " + strings.Join(args, " ")
    if err, ok := f.commands[key]; ok {
        return err
    }
    return fmt.Errorf("unexpected command: %s", key)
}

func TestCommandExecution(t *testing.T) {
    fakeRunner := &FakeCmdRunner{
        commands: map[string]error{
            "mvn clean": nil,
        },
    }
    
    err := ExecuteCommand(fakeRunner, "mvn", "clean")
    assert.NoError(t, err)
}
```

### 8. Testing Timestamp Comparison

```go
func TestTimestampCompare(t *testing.T) {
    tmpDir := t.TempDir()
    
    // Create older file
    oldFile := filepath.Join(tmpDir, "old.txt")
    err := os.WriteFile(oldFile, []byte("old"), 0644)
    require.NoError(t, err)
    
    // Wait to ensure different timestamp
    time.Sleep(10 * time.Millisecond)
    
    // Create newer file
    newFile := filepath.Join(tmpDir, "new.txt")
    err = os.WriteFile(newFile, []byte("new"), 0644)
    require.NoError(t, err)
    
    // Test
    isNewer, err := CompareTimestamps(oldFile, newFile)
    require.NoError(t, err)
    assert.False(t, isNewer) // newFile is newer
}
```

### 9. Using testdata Directory

```go
// testdata/java-project/pom.xml exists

func TestWithTestdata(t *testing.T) {
    projectRoot := filepath.Join("testdata", "java-project")
    
    ecosystems, err := DetectEcosystems(projectRoot, allConfigs)
    require.NoError(t, err)
    assert.Len(t, ecosystems, 1)
}
```

### 10. Parallel Tests

```go
func TestParallel(t *testing.T) {
    t.Parallel() // Enables parallel execution
    
    // Test code
}
```

## Helper Functions

### Setup Test Project

```go
func setupTestProject(t *testing.T, files map[string]string) string {
    tmpDir := t.TempDir()
    
    for path, content := range files {
        fullPath := filepath.Join(tmpDir, path)
        err := os.MkdirAll(filepath.Dir(fullPath), 0755)
        require.NoError(t, err)
        
        err = os.WriteFile(fullPath, []byte(content), 0644)
        require.NoError(t, err)
    }
    
    return tmpDir
}

// Usage
func TestExample(t *testing.T) {
    projectDir := setupTestProject(t, map[string]string{
        "pom.xml": `<?xml version="1.0"?>...`,
        "src/main/java/App.java": "public class App {}",
    })
    
    // Test with projectDir
}
```

## Test Commands Cheat Sheet

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./...

# Run with verbose output
go test -v ./...

# Run in parallel
go test -parallel 4 ./...

# Generate HTML coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. ./...

# Run integration tests
go test -tags=integration ./...

# Skip integration tests
go test -short ./...
```

## Common Assertions

```go
// Equality
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)

// Nil checks
assert.Nil(t, value)
assert.NotNil(t, value)

// Errors
assert.NoError(t, err)
assert.Error(t, err)
assert.ErrorIs(t, err, targetErr)

// Booleans
assert.True(t, condition)
assert.False(t, condition)

// Collections
assert.Len(t, slice, expectedLength)
assert.Contains(t, slice, item)
assert.Empty(t, slice)

// Files
assert.FileExists(t, path)
assert.DirExists(t, path)

// JSON
assert.JSONEq(t, expectedJSON, actualJSON)
```

## Best Practices Checklist

- [ ] Use `require` for setup that must succeed
- [ ] Use `assert` for validations
- [ ] Use `t.TempDir()` for temporary files
- [ ] Use table-driven tests for multiple cases
- [ ] Mock external dependencies
- [ ] Test error cases
- [ ] Use descriptive test names
- [ ] Keep tests independent
- [ ] Use `t.Parallel()` when safe
- [ ] Clean up resources (automatic with `t.TempDir()`)

