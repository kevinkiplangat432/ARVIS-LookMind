from fastapi import FastAPI
from app.core.config import settings
from app.api.routes.health import router as health_router


def create_app() -> FastAPI:
    app = FastAPI(
        title=settings.app_name,
        version=settings.app_version,
        docs_url="/docs",
        redoc_url="/redoc",
    )

    app.include_router(health_router, prefix="/api/v1", tags=["health"])

    return app


app = create_app()
