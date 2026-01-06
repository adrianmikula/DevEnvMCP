# Quick Start: Using Dev-Env Sentinel in Cursor

## Step 1: Build the Binary

```powershell
go build -o sentinel.exe ./cmd/sentinel
```

## Step 2: Configure Cursor

1. Open or create: `%APPDATA%\Cursor\mcp.json`

2. Add this configuration (update paths to match your setup):
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

3. **Restart Cursor** completely (close and reopen)

## Step 3: Test in Cursor

Open a chat in Cursor and try these commands:

### Check License Status
```
Check my license status
```

### Get Pro License Info
```
How do I get a Pro license?
```

### Verify Build Freshness (if you have a project)
```
Check if my build artifacts are up-to-date in [your project path]
```

## Troubleshooting

### Server Not Appearing
- Make sure you **restarted Cursor** after adding the config
- Check that the path to `sentinel.exe` is correct and absolute
- Verify the binary exists: `dir sentinel.exe`

### Config Errors
- Make sure `SENTINEL_CONFIG_DIR` points to your project root
- Verify `config/` directory exists with subdirectories
- Check that config files are valid YAML

### Tools Not Working
- Try `check_license_status` first (no arguments needed)
- Premium tools require a Pro license
- Some tools need a `project_root` argument

## What to Expect

When you ask Cursor to check your license status, it should:
1. Call the `check_license_status` MCP tool
2. Return information about your current tier (Free/Pro/Enterprise)
3. List available features

The MCP server runs in the background and communicates via stdio with Cursor.

