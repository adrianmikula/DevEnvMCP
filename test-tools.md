# Testing MCP Tools in Cursor

## What to Look For

When you ask me to use a tool, you should see:
1. A tool call indicator in the chat
2. The tool name being executed
3. A response from the tool

## Test Commands

Try these exact commands in our chat:

### Test 1: License Status
```
Check my license status
```

Expected: Should show Free tier with available features list.

### Test 2: Get Pro License
```
Get information about purchasing a Pro license
```

Expected: Should show Stripe and Apify payment links.

### Test 3: Verify Build (if you have a project)
```
Verify build freshness for E:\Source\DevEnvMCP
```

Expected: Should analyze the project and check build artifacts.

## Troubleshooting

If tools aren't being called:

1. **Check Cursor MCP Settings**:
   - Go to Cursor Settings
   - Look for MCP or Model Context Protocol section
   - Verify `dev-env-sentinel` is listed

2. **Check Cursor Logs**:
   - Look for errors related to MCP
   - Check if the server process starts

3. **Verify Configuration**:
   - Make sure `%APPDATA%\Cursor\mcp.json` exists
   - Verify the path to `sentinel.exe` is correct
   - Check that `SENTINEL_CONFIG_DIR` is set

4. **Test Binary Manually**:
   ```powershell
   # Should start and wait for input
   .\sentinel.exe
   ```

## What Should Happen

When you ask me to check your license status, I should:
1. Recognize the request
2. Call the `check_license_status` MCP tool
3. Return the license information

If this doesn't happen, there might be an issue with:
- MCP server registration in Cursor
- Path configuration
- Server startup

