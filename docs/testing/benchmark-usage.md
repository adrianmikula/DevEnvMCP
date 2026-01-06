# Running Benchmarks

## Quick Start

To run benchmarks, you need to:

1. **Navigate to the package directory** OR use the full package path
2. **Use the `-tags=benchmark` flag** to include benchmark tests

## Correct Commands

### From Project Root

```bash
# This will NOT work - Go tries to test root package
go test -bench=. -benchmem -tags=benchmark ./internal/mcp

# Instead, use one of these:
```

### Option 1: Navigate to Package Directory

```bash
cd internal/mcp
go test -bench=. -benchmem -tags=benchmark
```

### Option 2: Use Full Module Path

```bash
go test -bench=. -benchmem -tags=benchmark dev-env-sentinel/internal/mcp
```

### Option 3: Run from Project Root with Explicit Package

```bash
# From project root
go test -bench=. -benchmem -tags=benchmark -run=^$ dev-env-sentinel/internal/mcp
```

## Available Benchmarks

- `BenchmarkBootstrapTime` - MCP server bootstrap performance
- `BenchmarkEcosystemDetection_Small` - Small project detection
- `BenchmarkEcosystemDetection_Medium` - Medium project detection  
- `BenchmarkEcosystemDetection_Large` - Large project detection
- `BenchmarkBuildFreshness_Small` - Small project verification
- `BenchmarkBuildFreshness_Large` - Large project verification
- `BenchmarkMCPTool_VerifyBuildFreshness` - Full tool execution
- `BenchmarkConfigLoading` - Config file loading

## Example Output

```
goos: windows
goarch: amd64
pkg: dev-env-sentinel/internal/mcp
BenchmarkBootstrapTime-8                1000    1200000 ns/op    500000 B/op    5 allocs/op
BenchmarkEcosystemDetection_Small-8     1000    500000 ns/op    100000 B/op    2 allocs/op
```

## Performance Tests (No Build Tag Required)

For performance requirement tests (not benchmarks), run:

```bash
# From project root - this works fine
go test -v ./internal/mcp -run Test.*Performance
go test -v ./internal/mcp -run Test.*Time
```

These tests verify that operations complete within required time limits.

