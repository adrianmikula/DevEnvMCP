# Go Dependency Management Research

## Summary

**Go Modules** is the official, built-in dependency management system for Go. It comes with Go itself (since Go 1.11, default since Go 1.13) and requires **no separate installation**.

## Official Solution: Go Modules

### Status
- ✅ **Official**: Built into Go toolchain
- ✅ **Mature**: Introduced in Go 1.11 (2018), default since Go 1.13 (2019)
- ✅ **Trusted**: Used by entire Go ecosystem
- ✅ **Active**: Continuously maintained by Go team

### Key Features

1. **Semantic Versioning**: Uses semantic versioning (v1.2.3) for dependencies
2. **Decentralized**: Fetches directly from version control (Git, etc.)
3. **Reproducible Builds**: `go.mod` and `go.sum` ensure consistent builds
4. **Integrated**: Part of `go` command, no external tools needed
5. **Proxy Support**: Can use Go module proxy for faster downloads
6. **Workspace Support**: Go 1.18+ supports workspaces for multi-module projects

### Files

- **`go.mod`**: Module definition and dependencies
- **`go.sum`**: Checksums for dependency verification

## Deprecated/Unmaintained Tools

### ❌ Glide
- **Status**: Unmaintained
- **Replacement**: Go Modules
- **Note**: Was popular pre-Go 1.11, now deprecated

### ❌ Dep
- **Status**: Deprecated
- **Replacement**: Go Modules
- **Note**: Official experiment, superseded by Go Modules

### ❌ Other Tools
- **gopkg**: Older tool, not recommended
- **govendor**: Older tool, not recommended

## Go Modules Commands

### Initialize Module
```bash
go mod init <module-name>
```
Creates `go.mod` file. Module name typically matches repository path.

### Add Dependencies
```bash
# Add specific package
go get <package-path>

# Add with specific version
go get <package-path>@v1.2.3

# Add latest version
go get <package-path>@latest

# Update all dependencies
go get -u ./...
```

### Manage Dependencies
```bash
# Remove unused, add missing
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Vendor dependencies (optional)
go mod vendor
```

### View Dependencies
```bash
# List all dependencies
go list -m all

# Show dependency graph
go mod graph

# Why is a dependency included?
go mod why <package-path>
```

## Best Practices

### 1. Module Naming
```bash
# Use repository path
go mod init github.com/username/project-name

# Or use domain-based path
go mod init example.com/project-name
```

### 2. Version Pinning
- Go Modules automatically pins versions in `go.mod`
- Use semantic versioning: `v1.2.3`
- Use `go get -u` to update dependencies

### 3. Dependency Hygiene
```bash
# Regularly run to clean up
go mod tidy

# Check for updates
go list -m -u all
```

### 4. Security
```bash
# Check for vulnerabilities (Go 1.18+)
go list -json -m all | nancy sleuth

# Or use govulncheck (Go 1.18+)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### 5. Workspace Support (Go 1.18+)
For multi-module projects:
```bash
# Initialize workspace
go work init
go work use ./module1 ./module2
```

## Additional Tools (Optional)

### 1. **golangci-lint**
Static analysis tool (not a dependency manager, but useful)
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. **govulncheck**
Vulnerability checking (built into Go 1.18+)
```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### 3. **gopls**
Go language server (for IDE support)
```bash
go install golang.org/x/tools/gopls@latest
```

## Comparison: Go Modules vs Others

| Feature | Go Modules | Glide | Dep |
|--------|-----------|-------|-----|
| Official | ✅ | ❌ | ❌ |
| Built-in | ✅ | ❌ | ❌ |
| Maintained | ✅ | ❌ | ❌ |
| Semantic Versioning | ✅ | ✅ | ✅ |
| Reproducible Builds | ✅ | ✅ | ✅ |
| Proxy Support | ✅ | ❌ | ❌ |
| Workspace Support | ✅ | ❌ | ❌ |

## Migration from Old Tools

If you have an old project using Glide or Dep:

1. **Remove old tool files**: `glide.yaml`, `Gopkg.toml`, etc.
2. **Initialize Go module**: `go mod init <module-name>`
3. **Convert dependencies**: `go mod tidy` (auto-converts)
4. **Test**: `go build ./...`
5. **Remove vendor** (optional): Go Modules doesn't require vendor directory

## For Dev-Env Sentinel MCP

### Setup Steps

1. **Check Go version** (need 1.13+):
   ```bash
   go version
   ```

2. **Initialize module**:
   ```bash
   go mod init dev-env-sentinel
   ```

3. **Add dependencies** (as needed):
   ```bash
   go get github.com/stretchr/testify
   go get github.com/shirou/gopsutil/v3
   go get github.com/hellofresh/health-go/v4
   # ... etc
   ```

4. **Tidy up**:
   ```bash
   go mod tidy
   ```

## Conclusion

**No separate dependency manager needed!** Go Modules is built into Go and is the only tool you need. It's:
- Official and maintained
- Mature and stable
- Used by entire Go ecosystem
- No installation required

Just use `go mod` commands that come with Go.

