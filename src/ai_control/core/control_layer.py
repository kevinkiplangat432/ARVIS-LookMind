"""Core SDK functionality - ControlLayer main class"""

class ControlLayer:
    """Main SDK entry point for AI agent governance"""
    
    def __init__(self, api_key: str, org_id: str):
        self.api_key = api_key
        self.org_id = org_id
    
    def monitor(self, agent_id: str):
        """Decorator for monitoring agent functions"""
        def decorator(func):
            def wrapper(*args, **kwargs):
                # TODO: Implement monitoring logic
                return func(*args, **kwargs)
            return wrapper
        return decorator
