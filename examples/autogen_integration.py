"""Example: AutoGen integration with AI Control Layer"""

from src.ai_control import ControlLayer

def create_monitored_autogen_agent(control: ControlLayer, agent_id: str):
    """Create an AutoGen agent with AI Control Layer monitoring"""
    
    @control.monitor(agent_id=agent_id)
    def agent_function(message):
        # TODO: Implement AutoGen agent logic
        return f"Response to: {message}"
    
    return agent_function

# Usage example
if __name__ == "__main__":
    control = ControlLayer(api_key="your_key", org_id="your_org")
    agent = create_monitored_autogen_agent(control, "autogen-agent")
    
    result = agent("Hello, agent!")
    print(result)
