import time
import uuid
import httpx
from fastapi import APIRouter, HTTPException, Request
from pydantic import BaseModel
from app.core.config import settings

router = APIRouter()

ANTHROPIC_API_URL = "https://api.anthropic.com/v1/messages"
ANTHROPIC_API_VERSION = "2023-06-01"


class Message(BaseModel):
    role: str
    content: str


class ProxyRequest(BaseModel):
    model: str = "claude-haiku-4-5-20251001"
    max_tokens: int = 1024
    messages: list[Message]
    session_id: str | None = None


@router.post("/proxy/chat")
async def proxy_chat(payload: ProxyRequest):
    request_id = str(uuid.uuid4())
    session_id = payload.session_id or str(uuid.uuid4())
    started_at = time.time()

    headers = {
        "x-api-key": settings.anthropic_api_key,
        "anthropic-version": ANTHROPIC_API_VERSION,
        "content-type": "application/json",
    }

    body = {
        "model": payload.model,
        "max_tokens": payload.max_tokens,
        "messages": [m.model_dump() for m in payload.messages],
    }

    async with httpx.AsyncClient(timeout=60.0) as client:
        try:
            response = await client.post(
                ANTHROPIC_API_URL,
                headers=headers,
                json=body,
            )
            response.raise_for_status()
        except httpx.HTTPStatusError as e:
            raise HTTPException(
                status_code=e.response.status_code,
                detail=f"Anthropic API error: {e.response.text}",
            )
        except httpx.RequestError as e:
            raise HTTPException(
                status_code=503,
                detail=f"Could not reach Anthropic API: {str(e)}",
            )

    latency_ms = round((time.time() - started_at) * 1000)
    data = response.json()

    input_tokens = data.get("usage", {}).get("input_tokens", 0)
    output_tokens = data.get("usage", {}).get("output_tokens", 0)
    response_text = data.get("content", [{}])[0].get("text", "")

    print(f"""
--- Sentinel Intercepted ---
request_id  : {request_id}
session_id  : {session_id}
model       : {payload.model}
input_tokens: {input_tokens}
output_tokens:{output_tokens}
latency_ms  : {latency_ms}
----------------------------
    """)

    return {
        "request_id": request_id,
        "session_id": session_id,
        "model": payload.model,
        "response": response_text,
        "usage": {
            "input_tokens": input_tokens,
            "output_tokens": output_tokens,
        },
        "latency_ms": latency_ms,
    }
