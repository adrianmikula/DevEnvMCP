# Apify Billable Events - Quick Reference

## Billable Events List

Use this list to configure events in your Apify dashboard:

### Free Events (No Charge)
- `verify_build_freshness` - $0.00
- `check_infrastructure_parity` - $0.00
- `env_var_audit` - $0.00
- `check_license_status` - $0.00
- `get_pro_license` - $0.00

### Premium Events (Billable)
- `reconcile_environment` - **$0.05** ⭐ Most valuable
- `auto_fix` - **$0.05**
- `advanced_diagnostics` - **$0.03**
- `docker_orchestration` - **$0.10** (Enterprise)
- `custom_configs` - **$0.02** (Enterprise)

## Apify Dashboard Configuration

1. Go to Actor → Settings → Pricing
2. Enable "Pay-Per-Event (PPE)"
3. Add each billable event with its price
4. Set minimum: $0.00, maximum: $0.50

## How It Works

When an AI agent calls a premium tool:
1. Code detects it's a billable event
2. Logs event to stderr: `APIFY_EVENT:{"type":"reconcile_environment","price":0.05,...}`
3. Apify parses the log and bills the user
4. User is charged automatically

## Event Tracking

Events are automatically tracked when:
- Running in Apify environment (`APIFY_API_TOKEN` and `APIFY_ACTOR_ID` set)
- Tool is called by an AI agent
- Event has a price > $0.00

Free events are NOT logged (saves on logging overhead).

