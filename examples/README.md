# Examples

This directory contains integration examples for the AI Control Layer SDK.

## Available Examples

### LangChain Integration
`langchain_integration.py` - Shows how to integrate the SDK with LangChain agents using callbacks.

### AutoGen Integration
`autogen_integration.py` - Demonstrates wrapping AutoGen agents with monitoring.

## Running Examples

1. Ensure the SDK is installed: `pip install -e .`
2. Configure your `.env` file with API credentials
3. Run an example: `python examples/langchain_integration.py`

## Creating Custom Integrations

The SDK provides flexible hooks for integrating with any agent framework:

```python
from ai_control import ControlLayer

control = ControlLayer(api_key="...", org_id="...")

@control.monitor(agent_id="custom-agent")
def your_agent_function(input_data):
    # Your agent logic
    return result
```
