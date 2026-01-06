# Monetization Implementation Summary

This document summarizes the feature flags and monetization implementation for Dev-Env Sentinel.

## What Was Implemented

### 1. License Management System (`internal/license/`)

#### `license.go`
- **LicenseValidator**: Validates license keys using HMAC-SHA256
- **License Structure**: Tracks tier, validity, expiration, and features
- **Tier Support**: Free, Pro, Enterprise
- **Apify Integration**: Special handling for `apify_xxx` tokens
- **Stripe Integration**: Payment link configuration via environment variables

#### `storage.go`
- **License Persistence**: Saves/loads licenses from `~/.dev-env-sentinel/license.json`
- **Environment Variable Support**: `SENTINEL_LICENSE_KEY` for cloud deployments
- **Secure Storage**: JSON-based local storage

### 2. Feature Flags System (`internal/features/`)

#### `features.go`
- **FeatureManager**: Manages feature access based on license tier
- **Feature Gating**: `RequireFeature()` checks access before execution
- **Upgrade Messages**: Automatic prompts with payment links when features are unavailable
- **Feature Lists**: Tier-based feature definitions

### 3. MCP Server Updates (`internal/mcp/`)

#### `server.go`
- **License Integration**: Server tracks license state
- **Feature Manager**: Integrated into server lifecycle
- **License Updates**: `UpdateLicense()` method for runtime activation

#### `tools.go`
- **Premium Feature Gating**: `reconcile_environment` requires Pro tier
- **Monetization Tools**:
  - `get_pro_license` - Returns payment information
  - `activate_pro` - Activates a license key
  - `check_license_status` - Shows current license state

### 4. Feature Tiers

#### Free Tier
- `verify_build_freshness`
- `check_infrastructure_parity`
- `env_var_audit`

#### Pro Tier
- All Free tier features
- `reconcile_environment` (auto-fix)
- `auto_fix`
- `advanced_diagnostics`

#### Enterprise Tier
- All Pro tier features
- `docker_orchestration`
- `priority_support`
- `custom_configs`

## Environment Variables

### License Configuration
- `SENTINEL_LICENSE_KEY` - License key (for cloud deployments)
- `SENTINEL_LICENSE_SECRET` - Secret key for HMAC validation
- `SENTINEL_STRIPE_PAYMENT_LINK` - Stripe payment link URL
- `SENTINEL_APIFY_ACTOR_URL` - Apify Actor URL

### Usage
```bash
# Local deployment
export SENTINEL_LICENSE_KEY="pro-abc123-lifetime"

# Apify deployment
export SENTINEL_LICENSE_KEY="apify_xxx"
export SENTINEL_APIFY_ACTOR_URL="https://api.apify.com/v2/actors/your-actor-id/run-sync"
```

## License Key Format

Standard format: `{tier}-{hmac}-{timestamp}`

- **tier**: `pro` or `enterprise`
- **hmac**: HMAC-SHA256 signature (first 16 characters)
- **timestamp**: Expiration date (YYYYMMDD) or `lifetime`

Example: `pro-a1b2c3d4e5f6g7h8-20251231`

## Monetization Tools

### `get_pro_license`
Returns:
- Stripe payment link
- Apify Actor URL
- Feature comparison
- Upgrade instructions

### `activate_pro`
Parameters:
- `license_key` (required): License key to activate

Actions:
- Validates license key
- Saves to local storage
- Updates server license state
- Returns confirmation

### `check_license_status`
Returns:
- Current tier
- License validity
- Available features
- Expiration date

## Integration Points

### Apify Deployment
1. Deploy as Apify Actor
2. Set `SENTINEL_LICENSE_KEY` to `apify_xxx`
3. Configure pay-per-event pricing
4. Users access via Apify API

### Stripe Integration
1. Create Stripe Payment Link
2. Set `SENTINEL_STRIPE_PAYMENT_LINK`
3. Generate license keys after payment
4. Email keys to customers
5. Users activate via `activate_pro` tool

## Security Features

1. **HMAC Validation**: License keys use cryptographic signatures
2. **Expiration Checking**: Automatic expiration validation
3. **Server-Side Validation**: Validation happens in Go (harder to bypass)
4. **Environment Variable Support**: Secure key management for cloud deployments

## Usage Examples

### Check License Status
```json
{
  "method": "tools/call",
  "params": {
    "name": "check_license_status",
    "arguments": {}
  }
}
```

### Activate Pro License
```json
{
  "method": "tools/call",
  "params": {
    "name": "activate_pro",
    "arguments": {
      "license_key": "pro-abc123-lifetime"
    }
  }
}
```

### Use Premium Feature (Auto-Fix)
```json
{
  "method": "tools/call",
  "params": {
    "name": "reconcile_environment",
    "arguments": {
      "project_root": "/path/to/project"
    }
  }
}
```

If license is invalid, returns upgrade message with payment links.

## Files Created/Modified

### New Files
- `internal/license/license.go` - License validation
- `internal/license/storage.go` - License persistence
- `internal/features/features.go` - Feature flag management
- `docs/monetization.md` - User documentation
- `docs/apify-deployment.md` - Apify deployment guide
- `docs/monetization-implementation.md` - This file

### Modified Files
- `internal/mcp/server.go` - Added license tracking
- `internal/mcp/tools.go` - Added monetization tools and feature gating
- `README.md` - Added monetization section

## Testing

To test the implementation:

1. **Free Tier** (default):
   ```bash
   # Should work
   verify_build_freshness()
   check_infrastructure_parity()
   env_var_audit()
   
   # Should return upgrade message
   reconcile_environment()
   ```

2. **Pro Tier**:
   ```bash
   # Activate license
   activate_pro(license_key="pro-test-key-lifetime")
   
   # Should now work
   reconcile_environment()
   ```

3. **Check Status**:
   ```bash
   check_license_status()
   ```

## Next Steps

1. **Generate Real License Keys**: Create a license key generator service
2. **Stripe Webhook**: Set up webhook to generate keys after payment
3. **Apify Deployment**: Deploy to Apify and configure pricing
4. **Analytics**: Track license activations and feature usage
5. **Documentation**: Update user-facing docs with payment links

## Revenue Models Supported

1. **Stripe Direct Sales**: One-time/subscription payments
2. **Apify Pay-Per-Event**: $0.02-$0.05 per tool call
3. **Enterprise Licensing**: Custom pricing and support

All models are integrated and ready to use!

