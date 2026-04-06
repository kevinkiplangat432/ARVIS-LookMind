from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    app_name: str = "Sentinel"
    app_version: str = "0.1.0"
    app_env: str = "development"
    anthropic_api_key: str = ""

    class Config:
        env_file = ".env"


settings = Settings()
