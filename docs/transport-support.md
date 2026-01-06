# Transport Layer Support

Dev-Env Sentinel MCP server supports multiple transport layers for different deployment scenarios.

## Supported Transports

### 1. Stdio Transport (Local Use) ✅

**Status**: Fully implemented and tested

**Usage**: Default transport for local MCP clients like Cursor, Claude Desktop, etc.

**How it works**:
- Reads JSON-RPC messages from `stdin`
- Writes JSON-RPC responses to `stdout`
- Standard MCP protocol over standard input/output

**Configuration**: No configuration needed - this is the default.

**Example**:
```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "E:\\Source\\DevEnvMCP\\sentinel.exe",
      "args": []
    }
  }
}
```

### 2. SSE+HTTP Transport (Apify/Cloud) ✅

**Status**: Implemented

**Usage**: For cloud deployments, Apify Actors, and HTTP-based MCP clients.

**How it works**:
- HTTP POST endpoint `/message` for sending requests
- Server-Sent Events (SSE) endpoint `/sse` for receiving responses
- Health check endpoint `/health`

**Configuration**: Set environment variables:
```bash
export SENTINEL_TRANSPORT=sse
export SENTINEL_HTTP_PORT=8080
```

**Endpoints**:
- `POST /message` - Send MCP requests
- `GET /sse` - Server-Sent Events stream
- `GET /health` - Health check

**Example Request**:
```bash
curl -X POST http://localhost:8080/message \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/list",
    "params": {}
  }'
```

## Transport Detection

The server automatically detects which transport to use:

1. **Check `SENTINEL_TRANSPORT` environment variable**:
   - `sse` or `http` → Use SSE+HTTP transport
   - Not set → Use stdio transport

2. **Check `SENTINEL_HTTP_PORT` environment variable**:
   - If set → Use SSE+HTTP transport on that port
   - If not set and SSE mode → Default to port 8080

## Implementation Details

### Stdio Transport
- Uses `os.Stdin` and `os.Stdout`
- JSON decoder/encoder for message parsing
- Synchronous request/response pattern
- Best for local process communication

### SSE+HTTP Transport
- HTTP server with multiple endpoints
- SSE for streaming responses
- CORS headers included for web clients
- Best for cloud/serverless deployments

## Testing Transports

### Test Stdio Transport
```bash
# Send initialize request
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}' | ./sentinel.exe
```

### Test SSE+HTTP Transport
```bash
# Start server in SSE mode
export SENTINEL_TRANSPORT=sse
export SENTINEL_HTTP_PORT=8080
./sentinel.exe

# In another terminal, test health endpoint
curl http://localhost:8080/health

# Test message endpoint
curl -X POST http://localhost:8080/message \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
```

## Apify Integration

For Apify deployments:

1. **Set environment variables in Apify Actor**:
   ```
   SENTINEL_TRANSPORT=sse
   SENTINEL_HTTP_PORT=8080
   SENTINEL_LICENSE_KEY=apify_xxx
   ```

2. **Apify will call the HTTP endpoints**:
   - POST requests to `/message` endpoint
   - SSE connection to `/sse` for responses

3. **Health checks**:
   - Apify can use `/health` endpoint for monitoring

## Migration Between Transports

The server automatically detects the transport based on environment variables. No code changes needed:

- **Local development**: No env vars → Uses stdio
- **Apify deployment**: Set `SENTINEL_TRANSPORT=sse` → Uses SSE+HTTP

## Future Transport Support

Potential future transports:
- WebSocket transport (for real-time bidirectional communication)
- TCP transport (for network-based deployments)
- gRPC transport (for high-performance scenarios)

