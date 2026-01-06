# Monetization & Feature Flags

Dev-Env Sentinel supports a freemium model with feature flags and license-based access control.

## Pricing Tiers

### Free Tier
- ✅ Build freshness verification
- ✅ Infrastructure parity checks
- ✅ Environment variable auditing
- ❌ Auto-fix capabilities
- ❌ Advanced diagnostics

### Pro Tier
- ✅ All Free tier features
- ✅ Auto-fix environment issues (`reconcile_environment`)
- ✅ Advanced diagnostics
- ✅ Priority support

### Enterprise Tier
- ✅ All Pro tier features
- ✅ Docker orchestration
- ✅ Custom configurations
- ✅ Dedicated support

## License Activation

### Via Stripe Payment Link

1. Use the `get_pro_license` tool to get payment information
2. Purchase a license through the Stripe payment link
3. Receive your license key via email
4. Use the `activate_pro` tool with your license key:
   ```
   activate_pro(license_key="pro-abc123-lifetime")
   ```

### Via Apify (Pay-Per-Event)

For cloud deployments, you can use Apify's pay-per-event model:

1. Deploy the Sentinel as an Apify Actor
2. Set `SENTINEL_APIFY_ACTOR_URL` environment variable
3. Use Apify token format: `apify_xxx`
4. Each tool call is billed per event

## Environment Variables

### License Configuration

- `SENTINEL_LICENSE_KEY` - License key (for cloud deployments)
- `SENTINEL_LICENSE_SECRET` - Secret key for HMAC validation (server-side)
- `SENTINEL_STRIPE_PAYMENT_LINK` - Stripe payment link URL
- `SENTINEL_APIFY_ACTOR_URL` - Apify Actor URL for pay-per-event

### Apify Deployment

When deploying to Apify:

```bash
export SENTINEL_LICENSE_KEY="apify_xxx"
export SENTINEL_APIFY_ACTOR_URL="https://api.apify.com/v2/actors/your-actor-id/run-sync"
```

## License Key Format

License keys follow the format: `{tier}-{hmac}-{timestamp}`

- `tier`: `pro` or `enterprise`
- `hmac`: HMAC-SHA256 signature (first 16 chars)
- `timestamp`: Expiration date (YYYYMMDD) or `lifetime`

Example: `pro-a1b2c3d4e5f6g7h8-20251231`

## Feature Flags

Features are gated behind license tiers. The system automatically checks feature availability before executing premium tools.

### Checking License Status

Use the `check_license_status` tool to see:
- Current tier
- License validity
- Available features
- Expiration date (if applicable)

### Premium Feature Access

When a premium feature is requested without a valid license, the system returns:
- An error message explaining the feature is premium-only
- Upgrade instructions with payment links
- Information about available tiers

## Monetization Tools

### `get_pro_license`
Returns information about purchasing a Pro license, including:
- Stripe payment link
- Apify Actor URL
- Feature comparison

### `activate_pro`
Activates a Pro license with a license key:
- Validates the license key
- Saves to local storage (`~/.dev-env-sentinel/license.json`)
- Updates server license state
- Returns activation confirmation

### `check_license_status`
Shows current license status:
- Tier (free/pro/enterprise)
- Validity status
- Available features
- Expiration date

## Implementation Details

### License Storage

Licenses are stored in:
- **Local**: `~/.dev-env-sentinel/license.json`
- **Environment**: `SENTINEL_LICENSE_KEY` (for cloud deployments)

### License Validation

- Uses HMAC-SHA256 for key validation
- Checks expiration dates
- Validates tier and feature access
- Supports lifetime licenses

### Feature Gating

Premium features check license status before execution:
```go
if err := server.featureManager.RequireFeature("reconcile_environment"); err != nil {
    return upgradeMessage, err
}
```

## Apify Integration

### Deploying to Apify

1. Create an Apify Actor from the Go binary
2. Set up pay-per-event pricing ($0.02-$0.05 per tool call)
3. Configure environment variables
4. Users access via Apify API with token authentication

### Apify Token Format

Apify tokens are automatically recognized:
- Format: `apify_xxx`
- Validated by Apify infrastructure
- Grants Pro tier access

## Stripe Integration

### Setting Up Stripe Payment Links

1. Create a Stripe Payment Link
2. Configure automatic license key delivery via email
3. Set `SENTINEL_STRIPE_PAYMENT_LINK` environment variable
4. Users purchase and receive license keys automatically

### License Key Generation

For Stripe integration, you'll need to:
1. Generate license keys server-side after payment
2. Email keys to customers
3. Keys follow the format: `{tier}-{hmac}-{timestamp}`

## Security Considerations

- License keys use HMAC for validation (harder to bypass than plain text)
- Secret keys should be kept secure (use environment variables)
- License validation happens server-side
- Local storage is used for convenience, but cloud deployments use environment variables

## Revenue Models

### 1. Stripe (Direct Sales)
- One-time payments
- Subscription payments
- You keep ~97% after Stripe fees

### 2. Apify (Pay-Per-Event)
- $0.02-$0.05 per tool call
- Automatic billing
- Apify handles infrastructure

### 3. Enterprise Licensing
- Custom pricing
- Volume discounts
- Dedicated support

