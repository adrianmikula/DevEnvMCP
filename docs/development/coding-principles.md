# Coding Principles

## Core Principles

### 1. DRY (Don't Repeat Yourself)
- **Shared Utilities**: All common functionality in `internal/common`
- **Reusable Functions**: Path resolution, file operations, error handling
- **Single Source of Truth**: Configuration schema defined once

### 2. KISS (Keep It Simple, Stupid)
- **Focused Functions**: Each function does one thing well
- **Clear Names**: Descriptive function/variable names
- **Minimal Complexity**: Avoid over-engineering

### 3. Efficiency
- **Minimal Allocations**: Reuse objects when possible
- **Efficient File Operations**: Use `filepath.WalkDir` for directory traversal
- **Early Returns**: Fail fast, avoid deep nesting

### 4. Code Organization
- **Package Structure**: Logical grouping by functionality
- **Internal Packages**: Use `internal/` to prevent external imports
- **Clear Boundaries**: Each package has a clear responsibility

## Code Patterns

### Error Handling
```go
// Use structured errors from common package
if err != nil {
    return nil, &common.ErrNotFound{Resource: "file", Path: path}
}
```

### Path Operations
```go
// Always use common utilities
resolved, err := common.ResolvePath(path)
expanded := common.ExpandPattern(pattern)
```

### File Operations
```go
// Use common file utilities
if !common.FileExists(path) {
    return nil, &common.ErrNotFound{Resource: "file", Path: path}
}
```

## Anti-Patterns to Avoid

1. **Code Duplication**: If you write similar code twice, extract to `common/`
2. **Deep Nesting**: Use early returns, avoid >3 levels
3. **Large Functions**: Keep functions <50 lines when possible
4. **Magic Strings**: Use constants or config values
5. **Unnecessary Abstractions**: Don't abstract until you need to

## Code Review Checklist

- [ ] No code duplication (check `common/` first)
- [ ] Functions are focused and small
- [ ] Error handling is explicit
- [ ] Path operations use `common/` utilities
- [ ] No magic numbers or strings
- [ ] Clear, descriptive names
- [ ] Comments explain "why", not "what"

