import time
import uuid
import httpx
from fastapi import APIRouter, HTTPException, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from pydantic import BaseModel
from app.core.config import settings
from app.db.database import get_db
from app.db.models import InteractionLog

router = APIRouter()

ANTHROPIC_API_URL = "https://api.anthropic.com/v1/messages"
ANTHROPIC_VERSION = "2023-06-01"
GROQ_API_URL = "https://api.groq.com/openai/v1/chat/completions"


class Message(BaseModel):
    role: str
    content: str


class ProxyRequest(BaseModel):
    model: str | None = None
    max_tokens: int = 1024
    messages: list[Message]
    session_id: str | None = None
    provider: str | None = None


async def call_groq(messages: list, model: str, max_tokens: int) -> dict:
    headers = {
        "Authorization": f"Bearer {settings.groq_api_key}",
        "Content-Type": "application/json",
    }
    body = {"model": model, "max_tokens": max_tokens, "messages": messages}
    async with httpx.AsyncClient(timeout=60.0) as client:
        try:
            response = await client.post(GROQ_API_URL, headers=headers, json=body)
            response.raise_for_status()
        except httpx.HTTPStatusError as e:
            raise HTTPException(
                status_code=e.response.status_code,
                detail=f"Groq API error: {e.response.text}",
            )
        except httpx.RequestError as e:
            raise HTTPException(status_code=503, detail=f"Cannot reach Groq: {str(e)}")

    data = response.json()
    return {
        "text": data["choices"][0]["message"]["content"],
        "input_tokens": data["usage"]["prompt_tokens"],
        "output_tokens": data["usage"]["completion_tokens"],
    }


async def call_anthropic(messages: list, model: str, max_tokens: int) -> dict:
    headers = {
        "x-api-key": settings.anthropic_api_key,
        "anthropic-version": ANTHROPIC_VERSION,
        "content-type": "application/json",
    }
    body = {"model": model, "max_tokens": max_tokens, "messages": messages}
    async with httpx.AsyncClient(timeout=60.0) as client:
        try:
            response = await client.post(ANTHROPIC_API_URL, headers=headers, json=body)
            response.raise_for_status()
        except httpx.HTTPStatusError as e:
            raise HTTPException(
                status_code=e.response.status_code,
                detail=f"Anthropic API error: {e.response.text}",
            )
        except httpx.RequestError as e:
            raise HTTPException(status_code=503, detail=f"Cannot reach Anthropic: {str(e)}")

    data = response.json()
    return {
        "text": data["content"][0]["text"],
        "input_tokens": data["usage"]["input_tokens"],
        "output_tokens": data["usage"]["output_tokens"],
    }


PROVIDER_DEFAULTS = {
    "groq": "llama-3.1-8b-instant",
    "anthropic": "claude-haiku-4-5-20251001",
}


@router.post("/proxy/chat")
async def proxy_chat(payload: ProxyRequest, db: AsyncSession = Depends(get_db)):
    request_id = str(uuid.uuid4())
    session_id = payload.session_id or str(uuid.uuid4())
    provider = payload.provider or settings.default_provider
    model = payload.model or PROVIDER_DEFAULTS.get(provider, "llama-3.1-8b-instant")
    prompt = payload.messages[-1].content

    messages = [m.model_dump() for m in payload.messages]
    started_at = time.time()

    if provider == "groq":
        result = await call_groq(messages, model, payload.max_tokens)
    elif provider == "anthropic":
        result = await call_anthropic(messages, model, payload.max_tokens)
    else:
        raise HTTPException(status_code=400, detail=f"Unknown provider: {provider}")

    latency_ms = round((time.time() - started_at) * 1000)

    log = InteractionLog(
        request_id=request_id,
        session_id=session_id,
        provider=provider,
        model=model,
        prompt=prompt,
        response=result["text"],
        input_tokens=result["input_tokens"],
        output_tokens=result["output_tokens"],
        latency_ms=latency_ms,
    )
    db.add(log)
    await db.commit()

    print(f"""
--- Sentinel Intercepted ---
request_id   : {request_id}
session_id   : {session_id}
provider     : {provider}
model        : {model}
input_tokens : {result['input_tokens']}
output_tokens: {result['output_tokens']}
latency_ms   : {latency_ms}
logged to db : ✓
----------------------------
    """)

    return {
        "request_id": request_id,
        "session_id": session_id,
        "provider": provider,
        "model": model,
        "response": result["text"],
        "usage": {
            "input_tokens": result["input_tokens"],
            "output_tokens": result["output_tokens"],
        },
        "latency_ms": latency_ms,
    }
