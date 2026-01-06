# Apify Billable Events Configuration

This document lists all billable events for Dev-Env Sentinel MCP and how to configure them in your Apify dashboard.

## Event Categories

### Free Tier Events (No Charge)
These events are free and don't trigger billing:

| Event Type | Tool Name | Price | Description |
|------------|-----------|-------|-------------|
| `verify_build_freshness` | `verify_build_freshness` | $0.00 | Verify build artifact freshness |
| `check_infrastructure_parity` | `check_infrastructure_parity` | $0.00 | Check infrastructure service parity |
| `env_var_audit` | `env_var_audit` | $0.00 | Audit environment variables |
| `check_license_status` | `check_license_status` | $0.00 | Check license status |
| `get_pro_license` | `get_pro_license` | $0.00 | Get Pro license information |

### Premium Tier Events (Billable)
These events trigger billing when called:

| Event Type | Tool Name | Price | Description | Tier Required |
|------------|-----------|-------|-------------|---------------|
| `reconcile_environment` | `reconcile_environment` | **$0.05** | Auto-fix environment issues | Pro |
| `auto_fix` | (internal) | **$0.05** | Automatic issue resolution | Pro |
| `advanced_diagnostics` | (internal) | **$0.03** | Advanced diagnostic analysis | Pro |
| `docker_orchestration` | (future) | **$0.10** | Docker container orchestration | Enterprise |
| `custom_configs` | (future) | **$0.02** | Custom configuration management | Enterprise |

## Apify Dashboard Configuration

### Step 1: Configure Pay-Per-Event (PPE) Pricing

In your Apify Actor settings → Pricing:

1. **Enable Pay-Per-Event (PPE)**
2. **Set Base Price**: $0.00 (free tier tools)
3. **Configure Event-Based Pricing**:

```
Event: reconcile_environment
Price: $0.05
Description: Auto-fix environment issues

Event: auto_fix
Price: $0.05
Description: Automatic issue resolution

Event: advanced_diagnostics
Price: $0.03
Description: Advanced diagnostic analysis

Event: docker_orchestration
Price: $0.10
Description: Docker container orchestration

Event: custom_configs
Price: $0.02
Description: Custom configuration management
```

### Step 2: Set Pricing Limits

- **Minimum per run**: $0.00 (free tools don't charge)
- **Maximum per run**: $0.50 (safety limit for multiple premium operations)
- **Default price**: $0.00 (if no billable events occur)

### Step 3: Configure Event Tracking

The code automatically logs events in this format:
```
APIFY_EVENT:{"type":"reconcile_environment","tool_name":"reconcile_environment","price":0.05,"timestamp":"..."}
```

Apify will parse these logs and bill accordingly.

## Event Pricing Strategy

### Rationale

1. **Free Tier Tools** ($0.00):
   - Low compute cost
   - Build user base
   - Demonstrate value

2. **Premium Tools** ($0.03 - $0.05):
   - Medium compute cost
   - High value to users
   - Reasonable pricing for AI agents

3. **Enterprise Tools** ($0.10):
   - High compute cost (Docker operations)
   - Enterprise customers can afford it
   - Justifies premium pricing

## How Events Are Tracked

### Automatic Tracking

When a tool is called:
1. Code checks if event is billable (`IsBillableEvent()`)
2. If billable, logs event to stderr in JSON format
3. Apify parses logs and bills the user
4. Event includes metadata (user_id, project_root, etc.)

### Event Log Format

```json
{
  "type": "reconcile_environment",
  "tool_name": "reconcile_environment",
  "price": 0.05,
  "timestamp": "2024-01-01T00:00:00Z",
  "user_id": "apify_user_123",
  "project_root": "/tmp/project"
}
```

## Testing Event Tracking

### Test Free Event (No Charge)
```bash
# Call free tool
curl -X POST http://localhost:8080/message \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "check_license_status",
      "arguments": {}
    }
  }'

# Should NOT log APIFY_EVENT (free)
```

### Test Billable Event
```bash
# Call premium tool
curl -X POST http://localhost:8080/message \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "reconcile_environment",
      "arguments": {
        "project_root": "/tmp/test"
      }
    }
  }'

# Should log: APIFY_EVENT:{"type":"reconcile_environment",...}
```

## Revenue Projections

### Scenario 1: 1000 Free Tool Calls/Day
- Free events: 1000 × $0.00 = $0.00
- **Revenue**: $0.00/day

### Scenario 2: 1000 Free + 100 Premium Calls/Day
- Free events: 1000 × $0.00 = $0.00
- Premium events: 100 × $0.05 = $5.00
- **Revenue**: $5.00/day = $150/month
- **After Apify fees (15%)**: $127.50/month

### Scenario 3: Enterprise Usage (500 Premium + 50 Enterprise/Day)
- Premium events: 500 × $0.05 = $25.00
- Enterprise events: 50 × $0.10 = $5.00
- **Revenue**: $30.00/day = $900/month
- **After Apify fees (15%)**: $765/month

## Best Practices

1. **Log all billable events** - Even if tracking fails, log for audit
2. **Include metadata** - User ID, project root, etc. for analytics
3. **Set reasonable limits** - Prevent abuse with max execution time
4. **Monitor usage** - Track which events are most popular
5. **Adjust pricing** - Based on actual compute costs and demand

## Troubleshooting

### Events Not Being Billed

1. **Check event logging**:
   ```bash
   # Look for APIFY_EVENT logs in stderr
   grep "APIFY_EVENT" actor.log
   ```

2. **Verify Apify configuration**:
   - PPE pricing is enabled
   - Event names match exactly
   - Prices are set correctly

3. **Check environment variables**:
   ```bash
   echo $APIFY_API_TOKEN
   echo $APIFY_ACTOR_ID
   ```

### Incorrect Billing

1. **Review event logs** - Check what events were logged
2. **Verify pricing** - Ensure prices match dashboard
3. **Contact Apify support** - For billing disputes

## Event Metadata

Each event includes:
- `type`: Event type identifier
- `tool_name`: MCP tool name
- `price`: Billable amount
- `timestamp`: When event occurred
- `user_id`: Apify user identifier (if available)
- `project_root`: Project path (if applicable)

This metadata helps with:
- Analytics and usage tracking
- Debugging billing issues
- Understanding user behavior
- Optimizing pricing strategy

