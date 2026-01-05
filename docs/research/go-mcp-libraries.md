# Go MCP Server Libraries Research

This document compares available Go libraries and frameworks for building MCP (Model Context Protocol) servers to avoid reinventing the wheel.

## Overview

Several Go libraries exist for building MCP servers. This research compares their features, APIs, and suitability for the Dev-Env Sentinel MCP project.

## Library Comparison

### 1. Go MCP SDK (Official)

**Status**: Official SDK  
**URL**: https://www.mcpstack.org/sdks/golang/go-mcp-sdk  
**GitHub**: Likely `github.com/modelcontextprotocol/go-sdk` (verify)

**Features**:
- ✅ Official/authoritative implementation
- ✅ Native Go implementation of MCP protocol
- ✅ Type-safe APIs
- ✅ High-performance concurrent server/client
- ✅ Built-in middleware support
- ✅ Comprehensive testing utilities
- ✅ Full JSON-RPC 2.0 support
- ✅ Multiple transport layers

**Pros**:
- Official support from MCP maintainers
- Production-ready
- Well-maintained
- Full protocol compliance

**Cons**:
- May be newer (less community examples)
- Potentially more opinionated API

**Best For**: Production applications requiring official support

---

### 2. GoMCP

**Status**: Community library  
**URL**: https://gomcp.dev  
**GitHub**: Verify exact repository

**Features**:
- ✅ Transport-agnostic server logic
- ✅ Multiple transports: Stdio, SSE+HTTP, WebSocket, TCP
- ✅ Tool, resource, and prompt registration
- ✅ Argument parsing utilities
- ✅ Progress reporting
- ✅ Response creation helpers

**Pros**:
- Transport flexibility
- Good documentation (based on website)
- Idiomatic Go design

**Cons**:
- Community-maintained (less official)
- Need to verify GitHub repo and activity

**Best For**: Projects needing transport flexibility

---

### 3. mcp-golang (metoro-io)

**Status**: Community library  
**URL**: https://github.com/metoro-io/mcp-golang  
**GitHub**: `github.com/metoro-io/mcp-golang`

**Features**:
- ✅ Type safety through native Go structs
- ✅ Minimal boilerplate code
- ✅ Modular design
- ✅ Custom transports (stdio, HTTP)
- ✅ Automatic schema generation
- ✅ Dynamic tool/resource/prompt registration

**Pros**:
- Low boilerplate
- Type-safe
- Modular (use only what you need)
- Active development

**Cons**:
- Community-maintained
- Smaller ecosystem

**Best For**: Rapid prototyping and projects valuing type safety

---

### 4. mcp-go (FlexInfer)

**Status**: Community library  
**URL**: https://pkg.go.dev/gitlab.flexinfer.ai/libs/mcp-go  
**GitHub**: `gitlab.flexinfer.ai/libs/mcp-go`

**Features**:
- ✅ Complete MCP protocol implementation
- ✅ Full JSON-RPC 2.0 support
- ✅ Multiple transports (Stdio, WebSocket)
- ✅ Simple API for tool registration
- ✅ Production-ready
- ✅ Extensively tested

**Pros**:
- Production-proven
- Simple API
- Good test coverage
- Complete protocol support

**Cons**:
- Hosted on GitLab (not GitHub)
- May have different community

**Best For**: Projects needing proven, production-ready solution

---

### 5. MCP-Go

**Status**: Community library  
**URL**: https://mcp-go.dev  
**GitHub**: Verify exact repository

**Features**:
- ✅ Define tools/resources/prompts as Go functions
- ✅ Automatic schema generation
- ✅ Multiple transport layers (Stdio, HTTP, SSE, Gin)
- ✅ Low boilerplate
- ✅ Type-safe

**Pros**:
- Function-based API (clean)
- Automatic schema generation
- Multiple transport options

**Cons**:
- Need to verify repository and maintenance

**Best For**: Projects wanting function-based API

---

## Feature Matrix

| Feature | Go MCP SDK | GoMCP | mcp-golang | mcp-go | MCP-Go |
|---------|-----------|-------|------------|--------|--------|
| Official Support | ✅ | ❌ | ❌ | ❌ | ❌ |
| Stdio Transport | ✅ | ✅ | ✅ | ✅ | ✅ |
| HTTP Transport | ✅ | ✅ | ✅ | ❌ | ✅ |
| WebSocket Transport | ✅ | ✅ | ❌ | ✅ | ❌ |
| Type Safety | ✅ | ✅ | ✅ | ✅ | ✅ |
| Low Boilerplate | ✅ | ✅ | ✅ | ✅ | ✅ |
| Schema Generation | ✅ | ✅ | ✅ | ❌ | ✅ |
| Middleware | ✅ | ❌ | ❌ | ❌ | ❌ |
| Testing Utilities | ✅ | ❌ | ❌ | ✅ | ❌ |
| Production Ready | ✅ | ✅ | ✅ | ✅ | ✅ |

## Recommendation for Dev-Env Sentinel MCP

### Primary Recommendation: **Go MCP SDK** (if available and stable)

**Rationale**:
1. **Official Support**: Official SDK means long-term maintenance and protocol compliance
2. **Production Ready**: Built for production use
3. **Full Feature Set**: All transports, middleware, testing utilities
4. **Future-Proof**: Will stay aligned with MCP protocol evolution

**Action**: Verify GitHub repository, check recent activity, and review API

### Alternative Recommendation: **mcp-golang** (metoro-io)

**Rationale**:
1. **Low Boilerplate**: Fits our configuration-driven architecture
2. **Type Safety**: Important for our ecosystem config system
3. **Modular**: Can use only what we need
4. **Active Development**: Appears to be actively maintained

**Action**: Review GitHub repository, check examples, verify API matches our needs

### Fallback: **mcp-go** (FlexInfer)

**Rationale**:
1. **Production Proven**: Used in multiple MCP servers
2. **Simple API**: Easy to integrate
3. **Complete**: Full protocol support

**Action**: Review if other options don't work out

## Evaluation Criteria

When evaluating libraries, consider:

1. **Protocol Compliance**: Does it fully implement MCP spec?
2. **API Design**: Does it fit our architecture (config-driven)?
3. **Transport Support**: Do we need only Stdio, or also HTTP/WebSocket?
4. **Maintenance**: Is it actively maintained?
5. **Documentation**: Is there good documentation and examples?
6. **Community**: Are there examples and community support?
7. **Testing**: Does it include testing utilities?
8. **Performance**: Is it performant enough for our use case?

## Next Steps

1. **Verify Repositories**: Check actual GitHub/GitLab repositories for each library
2. **Review APIs**: Look at code examples and API documentation
3. **Check Activity**: Review recent commits, issues, and releases
4. **Test Integration**: Create a small proof-of-concept with top 2-3 candidates
5. **Make Decision**: Choose based on fit for our architecture

## Example Integration (Conceptual)

### Using Go MCP SDK (if available)

```go
import (
    "context"
    mcp "github.com/modelcontextprotocol/go-sdk/server"
)

func main() {
    server := mcp.NewServer("dev-env-sentinel")
    
    // Register tools
    server.RegisterTool("verify_build_freshness", verifyBuildFreshness)
    server.RegisterTool("check_infrastructure_parity", checkInfrastructureParity)
    server.RegisterTool("env_var_audit", envVarAudit)
    server.RegisterTool("reconcile_environment", reconcileEnvironment)
    
    // Start stdio transport
    transport := mcp.NewStdioTransport()
    server.Serve(transport)
}
```

### Using mcp-golang (alternative)

```go
import (
    "github.com/metoro-io/mcp-golang/server"
)

func main() {
    srv := server.NewServer()
    
    srv.RegisterTool("verify_build_freshness", &VerifyBuildFreshnessTool{})
    srv.RegisterTool("check_infrastructure_parity", &CheckInfrastructureParityTool{})
    
    // Use stdio transport
    transport := server.NewStdioTransport(srv)
    transport.Serve()
}
```

## Additional Resources

- **Awesome MCP DevTools**: https://github.com/punkpeye/awesome-mcp-devtools
  - Curated list of MCP tools, SDKs, and libraries
  - Good starting point for discovering more options

- **MCP Specification**: https://modelcontextprotocol.io
  - Official protocol specification
  - Reference for evaluating library compliance

## Decision Framework

1. **If Go MCP SDK exists and is stable**: Use it (official support)
2. **If Go MCP SDK doesn't exist or is unstable**: Evaluate mcp-golang and mcp-go
3. **If we need custom features**: Consider forking or contributing to chosen library
4. **If all libraries are insufficient**: Build minimal MCP implementation (JSON-RPC over stdio)

## Notes

- **Stdio Transport**: Most important for MCP servers (Claude Desktop uses stdio)
- **HTTP/WebSocket**: Nice to have for future web integration, but not required for MVP
- **Schema Generation**: Helpful but not critical (we can define schemas manually)
- **Middleware**: Useful for logging, error handling, but can be added later

## Action Items

- [ ] Research actual GitHub repositories for each library
- [ ] Review API documentation and examples
- [ ] Check recent activity and maintenance status
- [ ] Create proof-of-concept with top candidate
- [ ] Make final decision based on fit for our architecture
- [ ] Document chosen library in architecture docs

