"""Configuration management"""

from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    # Application
    app_name: str = "ai-control-layer"
    app_version: str = "0.1.0"
    environment: str = "development"
    debug: bool = True
    
    # API
    api_host: str = "0.0.0.0"
    api_port: int = 8000
    api_key: Optional[str] = None
    
    # Database
    database_url: str = "postgresql+asyncpg://user:password@localhost:5432/ai_control"
    database_pool_size: int = 20
    database_max_overflow: int = 10
    
    # Redis
    redis_url: str = "redis://localhost:6379/0"
    redis_max_connections: int = 50
    
    # Security
    secret_key: str = "change-me-in-production"
    jwt_algorithm: str = "HS256"
    jwt_expiration_hours: int = 24
    
    # Risk Engine
    risk_threshold_high: float = 0.8
    risk_threshold_medium: float = 0.5
    risk_threshold_low: float = 0.3
    
    class Config:
        env_file = ".env"

settings = Settings()
