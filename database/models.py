"""Database models using SQLAlchemy"""

from sqlalchemy import Column, String, Integer, DateTime, JSON
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

Base = declarative_base()

class AgentLog(Base):
    __tablename__ = "agent_logs"
    
    id = Column(Integer, primary_key=True)
    agent_id = Column(String, nullable=False)
    event_type = Column(String, nullable=False)
    payload = Column(JSON)
    timestamp = Column(DateTime, default=datetime.utcnow)

class RiskScore(Base):
    __tablename__ = "risk_scores"
    
    id = Column(Integer, primary_key=True)
    agent_id = Column(String, nullable=False)
    score = Column(Integer)
    metadata = Column(JSON)
    timestamp = Column(DateTime, default=datetime.utcnow)
