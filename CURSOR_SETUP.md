# Cursor MCP Setup Guide

This guide explains how to set up and use Dev-Env Sentinel as an MCP server in Cursor.

## Quick Setup

### Option 1: Using the Built Binary (Recommended for Development)

1. **Build the binary** (if not already built):
   ```bash
   go build -o sentinel.exe ./cmd/sentinel
   ```

2. **Configure Cursor MCP**:
   
   On Windows, edit or create: `%APPDATA%\Cursor\mcp.json`
   
   On Linux/Mac, edit or create: `~/.cursor/mcp.json`
   
   Copy the configuration from `cursor-mcp.json.example` and update the paths:
   ```json
   {
     "mcpServers": {
       "dev-env-sentinel": {
         "command": "E:\\Source\\DevEnvMCP\\sentinel.exe",
         "args": [],
         "env": {
           "SENTINEL_CONFIG_DIR": "E:\\Source\\DevEnvMCP"
         }
       }
     }
   }
   ```
   
   **Important**: 
   - Update the `command` path to your actual `sentinel.exe` location
   - Update `SENTINEL_CONFIG_DIR` to your project root directory
   - Use absolute paths on Windows (with double backslashes `\\`)

3. **Restart Cursor** to load the MCP server.

### Option 2: Using npm/npx (For Published Package)

Once published to npm:

```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "npx",
      "args": ["-y", "dev-env-sentinel"],
      "env": {}
    }
  }
}
```

## Configuration Details

### Environment Variables

- `SENTINEL_CONFIG_DIR` - Explicit path to the config directory (optional)
  - If not set, the server will try to find configs relative to the executable
  - Falls back to current working directory

- `SENTINEL_LICENSE_KEY` - License key for Pro features (optional)
  - Can be set here or activated via `activate_pro` tool

### Path Resolution

The server tries to find configs in this order:
1. `SENTINEL_CONFIG_DIR` environment variable
2. `config/` directory relative to executable
3. `../config/` relative to executable (for npm package structure)
4. Current working directory

## Testing the Setup

1. **Restart Cursor** after adding the MCP configuration

2. **Check if the server is loaded**:
   - Open Cursor's MCP panel (if available)
   - Or try using one of the tools in a chat

3. **Test a simple tool**:
   ```
   Check my license status
   ```
   
   This should call the `check_license_status` tool.

4. **Test config discovery**:
   ```
   Verify build freshness for the current project
   ```
   
   This should call `verify_build_freshness` with the current directory.

## Available Tools

### Free Tier Tools
- `verify_build_freshness` - Check if build artifacts are up-to-date
- `check_infrastructure_parity` - Verify infrastructure services
- `env_var_audit` - Audit environment variables

### Premium Tools (Require Pro License)
- `reconcile_environment` - Auto-fix environment issues

### Monetization Tools
- `get_pro_license` - Get information about purchasing Pro
- `activate_pro` - Activate a Pro license key
- `check_license_status` - Check current license status

## Troubleshooting

### Server Not Loading

1. **Check the binary path**:
   - Make sure the path in `mcp.json` is correct
   - Use absolute paths on Windows
   - Check that the binary exists and is executable

2. **Check Cursor logs**:
   - Look for MCP-related errors in Cursor's output
   - Check if the server process starts

3. **Test the binary manually**:
   ```bash
   # Should start MCP server (reads from stdin)
   ./sentinel.exe
   ```

### Config Files Not Found

1. **Set explicit config directory**:
   ```json
   "env": {
     "SENTINEL_CONFIG_DIR": "E:\\Source\\DevEnvMCP"
   }
   ```

2. **Check config directory structure**:
   ```
   config/
   ├── languages/
   │   ├── java.yaml
   │   ├── python.yaml
   │   └── javascript.yaml
   ├── tools/
   │   ├── java/
   │   ├── python/
   │   └── javascript/
   └── infrastructure/
       └── docker/
   ```

3. **Verify config files exist**:
   ```bash
   dir config\languages
   dir config\tools
   ```

### Tools Not Working

1. **Check license status**:
   - Use `check_license_status` tool
   - Premium tools require Pro license

2. **Check project root**:
   - Tools need a `project_root` argument
   - Make sure you're in a project directory

3. **Check tool arguments**:
   - Each tool has specific required arguments
   - Check tool descriptions for details

## Example Usage in Cursor

### Check License Status
```
User: Check my license status
AI: [Calls check_license_status tool]
```

### Verify Build Freshness
```
User: Check if my build artifacts are up-to-date in E:\Source\MyProject
AI: [Calls verify_build_freshness with project_root="E:\\Source\\MyProject"]
```

### Get Pro License Info
```
User: How do I get a Pro license?
AI: [Calls get_pro_license tool]
```

### Activate Pro License
```
User: Activate my Pro license with key pro-abc123-lifetime
AI: [Calls activate_pro with license_key="pro-abc123-lifetime"]
```

## Next Steps

1. **Test all tools** to ensure they work correctly
2. **Set up license keys** if you want to test Pro features
3. **Configure Stripe/Apify** links for actual monetization
4. **Publish to npm** for easier distribution

## Notes

- The MCP server uses stdio for communication (standard input/output)
- All errors are logged to stderr
- The server runs in the background when Cursor is active
- License keys are stored in `~/.dev-env-sentinel/license.json`

