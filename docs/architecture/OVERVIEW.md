# Architecture Overview

## System Components

### 1. SDK Layer (src/ai_control/)
- **core/**: Main ControlLayer class and SDK entry point
- **interceptors/**: Agent action interception and wrapping
- **events/**: Event streaming and queue management
- **risk/**: Risk evaluation algorithms and scoring
- **policy/**: Policy enforcement and RBAC
- **telemetry/**: Observability metrics and logging

### 2. API Layer (api/)
- FastAPI-based REST API
- WebSocket support for real-time events
- Authentication and authorization
- Rate limiting and request validation

### 3. Database Layer (database/)
- PostgreSQL for persistent storage
- SQLAlchemy ORM with async support
- Alembic for schema migrations
- Models for logs, events, risk scores, policies

### 4. Event Streaming
- Redis for lightweight event queuing
- Kafka support for high-throughput scenarios
- Async event processing pipelines

## Data Flow

1. Agent action → SDK interceptor
2. Event creation → Event queue (Redis/Kafka)
3. Risk evaluation → Risk score calculation
4. Policy check → Allow/Block/Request approval
5. Logging → Database persistence
6. Telemetry → Metrics export

## Scalability Considerations

- Async-first design for high concurrency
- Horizontal scaling via stateless API
- Event queue for decoupled processing
- Database connection pooling
- Caching layer for policies and configurations
