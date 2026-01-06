# Apify Deployment Guide

This guide explains how to deploy Dev-Env Sentinel as an Apify Actor for pay-per-event monetization.

## Overview

Apify provides a marketplace for "Actors" (serverless functions) with built-in billing. This allows you to monetize the Sentinel MCP without managing infrastructure or payment processing.

## Benefits

- ✅ No infrastructure management
- ✅ Automatic billing per execution
- ✅ Built-in API endpoints
- ✅ Apify handles scaling
- ✅ Fast payout (monthly)

## Prerequisites

1. Apify account (https://apify.com)
2. Go binary compiled for Linux
3. Environment variables configured

## Step 1: Create Apify Actor

1. Go to Apify Console → Actors → Create New
2. Choose "Docker" as the base image
3. Name it: `dev-env-sentinel`

## Step 2: Create Dockerfile

Create a `Dockerfile` in your project:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sentinel ./cmd/sentinel

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sentinel .
COPY --from=builder /app/config ./config
ENV PATH="/root:$PATH"
CMD ["./sentinel"]
```

## Step 3: Configure Environment Variables

In Apify Actor settings, add:

- `SENTINEL_LICENSE_KEY` - Set to `apify_xxx` (Apify handles this)
- `SENTINEL_APIFY_ACTOR_URL` - Your Actor URL
- `SENTINEL_LICENSE_SECRET` - Your secret key for license validation

## Step 4: Set Pricing

In Apify Actor settings → Pricing:

- **Pay-Per-Event (PPE)**: $0.02 - $0.05 per execution
- **Minimum**: $0.01 per run
- **Maximum**: $1.00 per run (for safety)

## Step 5: Deploy

1. Push your code to GitHub
2. Connect Apify to your repository
3. Set build command: `docker build -t dev-env-sentinel .`
4. Deploy

## Step 6: Test

Test your Actor:

```bash
curl -X POST https://api.apify.com/v2/actors/your-actor-id/run-sync \
  -H "Authorization: Bearer YOUR_APIFY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "method": "tools/call",
      "params": {
        "name": "verify_build_freshness",
        "arguments": {
          "project_root": "/tmp/test"
        }
      }
    }
  }'
```

## Step 7: Publish to Apify Store

1. Go to Actor → Settings → Store
2. Add description and tags
3. Set pricing
4. Publish

## Usage from MCP Clients

Users can configure the Apify Actor URL in their MCP settings:

```json
{
  "mcpServers": {
    "dev-env-sentinel": {
      "command": "curl",
      "args": [
        "-X", "POST",
        "https://api.apify.com/v2/actors/your-actor-id/run-sync",
        "-H", "Authorization: Bearer YOUR_APIFY_TOKEN",
        "-H", "Content-Type: application/json",
        "-d", "@-"
      ]
    }
  }
}
```

Or use a wrapper script that converts MCP stdio to HTTP requests.

## Revenue Model

- **Per Execution**: $0.02 - $0.05
- **Apify Fee**: ~10-15% (varies by plan)
- **Your Revenue**: ~85-90% of charges

Example:
- 1000 executions/day × $0.03 = $30/day
- Monthly: ~$900
- After Apify fees: ~$765/month

## Best Practices

1. **Set reasonable limits** - Prevent abuse with max execution time
2. **Monitor usage** - Track costs and optimize
3. **Cache results** - Reduce redundant executions
4. **Error handling** - Graceful failures don't charge users
5. **Documentation** - Clear usage instructions

## Troubleshooting

### Actor fails to start
- Check Dockerfile syntax
- Verify Go binary is compiled for Linux
- Check environment variables

### Billing issues
- Verify pricing is set correctly
- Check Apify account balance
- Review execution logs

### Performance
- Optimize Go binary size
- Use Alpine Linux base image
- Enable caching for faster cold starts

## Alternative: Apify SDK

For more control, use the Apify Go SDK:

```go
import "github.com/apify/apify-sdk-go"

// Initialize Apify client
client := apify.NewClient(os.Getenv("APIFY_TOKEN"))

// Handle requests
// ...
```

This allows more advanced features like:
- Custom request handling
- Better error messages
- Advanced logging
- Metrics collection

