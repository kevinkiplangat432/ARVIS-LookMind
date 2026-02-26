"""Pydantic schemas for data validation"""

from pydantic import BaseModel
from datetime import datetime
from typing import Optional, Dict, Any

class AgentLogSchema(BaseModel):
    agent_id: str
    event_type: str
    payload: Dict[str, Any]
    timestamp: Optional[datetime] = None

class RiskScoreSchema(BaseModel):
    agent_id: str
    score: float
    metadata: Dict[str, Any]
    timestamp: Optional[datetime] = None

class EventSchema(BaseModel):
    event_id: str
    event_type: str
    agent_id: str
    data: Dict[str, Any]
    timestamp: datetime
