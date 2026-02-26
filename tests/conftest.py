"""Pytest configuration and fixtures"""

import pytest
import asyncio
from typing import Generator

@pytest.fixture(scope="session")
def event_loop() -> Generator:
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()

@pytest.fixture
def mock_control_layer():
    """Mock ControlLayer instance for testing"""
    from src.ai_control import ControlLayer
    return ControlLayer(api_key="test_key", org_id="test_org")
