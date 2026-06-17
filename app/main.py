from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.core.config import settings
from app.api.routes.health import router as health_router
from app.api.routes.proxy import router as proxy_router
from app.db.database import engine, Base
from app.db import models  # noqa: F401 — ensures models are registered


@asynccontextmanager
async def lifespan(app: FastAPI):
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield


def create_app() -> FastAPI:
    app = FastAPI(
        title=settings.app_name,
        version=settings.app_version,
        docs_url="/docs",
        redoc_url="/redoc",
        lifespan=lifespan,
    )

    app.include_router(health_router, prefix="/api/v1", tags=["health"])
    app.include_router(proxy_router, prefix="/api/v1", tags=["proxy"])

    return app


app = create_app()
