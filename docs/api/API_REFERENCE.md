# API Reference

## Authentication

All API requests require authentication via API key in the header:
```
Authorization: Bearer <api_key>
```

## Endpoints

### Health Check
```
GET /health
```

### Agent Logs
```
POST /api/v1/logs
GET /api/v1/logs/{agent_id}
```

### Risk Scores
```
POST /api/v1/risk/evaluate
GET /api/v1/risk/{agent_id}
```

### Policies
```
GET /api/v1/policies
POST /api/v1/policies
PUT /api/v1/policies/{policy_id}
DELETE /api/v1/policies/{policy_id}
```

### Events
```
GET /api/v1/events
WebSocket /ws/events
```

## SDK Methods

### ControlLayer

```python
control = ControlLayer(api_key="...", org_id="...")
```

### Monitor Decorator

```python
@control.monitor(agent_id="my-agent")
def agent_function(input_data):
    return process(input_data)
```

## Event Types

- `PROMPT_SENT`
- `RESPONSE_RECEIVED`
- `TOOL_CALLED`
- `TOOL_RESULT`
- `AGENT_LOOP_ITERATION`
- `RISK_EVALUATED`
- `POLICY_CHECKED`
