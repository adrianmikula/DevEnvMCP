# PowerShell script to test the MCP server manually
# This simulates what Cursor would send to the server

Write-Host "Testing MCP Server..." -ForegroundColor Green
Write-Host ""

# Test 1: Initialize request
Write-Host "Test 1: Initialize Request" -ForegroundColor Yellow
$initRequest = @{
    jsonrpc = "2.0"
    id = 1
    method = "initialize"
    params = @{
        protocolVersion = "2024-11-05"
        capabilities = @{}
        clientInfo = @{
            name = "test-client"
            version = "1.0.0"
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Sending initialize request..."
$initRequest | .\sentinel.exe

Write-Host ""
Write-Host "Test 2: Tools List Request" -ForegroundColor Yellow
$toolsListRequest = @{
    jsonrpc = "2.0"
    id = 2
    method = "tools/list"
    params = @{}
} | ConvertTo-Json -Depth 10

Write-Host "Sending tools/list request..."
$toolsListRequest | .\sentinel.exe

Write-Host ""
Write-Host "Test 3: Check License Status" -ForegroundColor Yellow
$checkLicenseRequest = @{
    jsonrpc = "2.0"
    id = 3
    method = "tools/call"
    params = @{
        name = "check_license_status"
        arguments = @{}
    }
} | ConvertTo-Json -Depth 10

Write-Host "Sending check_license_status request..."
$checkLicenseRequest | .\sentinel.exe

Write-Host ""
Write-Host "Testing complete!" -ForegroundColor Green
Write-Host "Note: The server reads from stdin, so you'll need to pipe JSON to it."
Write-Host "Example: Get-Content test-request.json | .\sentinel.exe"

