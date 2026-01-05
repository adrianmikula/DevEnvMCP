# MCP Library Decision Guide

Quick reference for choosing the Go MCP library for Dev-Env Sentinel MCP.

## Quick Decision Tree

```
Is Go MCP SDK (official) available and stable?
├─ YES → Use Go MCP SDK (official support, future-proof)
└─ NO → Continue evaluation
    │
    Do we need minimal boilerplate + type safety?
    ├─ YES → Use mcp-golang (metoro-io)
    └─ NO → Continue evaluation
        │
        Do we need production-proven solution?
        ├─ YES → Use mcp-go (FlexInfer)
        └─ NO → Use GoMCP (transport flexibility)
```

## Top 3 Candidates

### 1. Go MCP SDK (Official) - **PREFERRED**

**Why**: Official support, full features, future-proof

**Check**:
- [ ] Repository exists: `github.com/modelcontextprotocol/go-sdk`
- [ ] Recent commits (within last 3 months)
- [ ] Documentation available
- [ ] API matches our needs

**If available**: Use this

---

### 2. mcp-golang (metoro-io) - **ALTERNATIVE**

**Why**: Low boilerplate, type-safe, modular

**Check**:
- [ ] Repository: `github.com/metoro-io/mcp-golang`
- [ ] Recent activity
- [ ] API examples available
- [ ] Stdio transport works

**If Go MCP SDK unavailable**: Use this

---

### 3. mcp-go (FlexInfer) - **FALLBACK**

**Why**: Production-proven, simple API

**Check**:
- [ ] Repository: `gitlab.flexinfer.ai/libs/mcp-go`
- [ ] Recent activity
- [ ] Examples available

**If others don't work**: Use this

## Minimum Requirements

For Dev-Env Sentinel MCP, we need:

- ✅ **Stdio Transport** (required for Claude Desktop)
- ✅ **Tool Registration** (for our 4 tools)
- ✅ **JSON-RPC 2.0** (MCP protocol requirement)
- ✅ **Type Safety** (for our config-driven architecture)
- ✅ **Low Boilerplate** (to keep code clean)

**Nice to Have**:
- HTTP/WebSocket transport (future)
- Schema generation (convenience)
- Middleware (logging, error handling)

## Verification Checklist

Before choosing a library:

- [ ] Repository exists and is accessible
- [ ] Recent commits (active maintenance)
- [ ] Documentation available
- [ ] Examples or tutorials exist
- [ ] API matches our architecture needs
- [ ] License is compatible (MIT, Apache 2.0, etc.)
- [ ] No major open issues blocking usage
- [ ] Community support (issues, discussions)

## Integration Test

Create a minimal test to verify library works:

```go
// test_mcp_library.go
package main

import (
    "context"
    // Import chosen library
)

func main() {
    // 1. Create server
    // 2. Register one tool
    // 3. Start stdio transport
    // 4. Verify it works
}
```

## Decision Timeline

1. **Week 1**: Research and verify repositories
2. **Week 1**: Review APIs and documentation
3. **Week 1**: Create proof-of-concept
4. **Week 1**: Make decision
5. **Week 2**: Integrate chosen library

## Notes

- **Stdio is critical**: Claude Desktop uses stdio transport
- **We can switch later**: If chosen library doesn't work out, we can refactor
- **Minimal implementation**: If all libraries fail, we can implement minimal MCP (JSON-RPC over stdio is straightforward)

