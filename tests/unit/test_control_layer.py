"""Unit tests for ControlLayer"""

import pytest
from src.ai_control import ControlLayer

def test_control_layer_initialization():
    control = ControlLayer(api_key="test", org_id="org")
    assert control.api_key == "test"
    assert control.org_id == "org"

def test_monitor_decorator(mock_control_layer):
    @mock_control_layer.monitor(agent_id="test-agent")
    def sample_function(x):
        return x * 2
    
    result = sample_function(5)
    assert result == 10
