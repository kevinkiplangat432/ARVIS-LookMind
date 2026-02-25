# AI Control Layer - Technical Documentation

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Problem Definition and Market Context](#problem-definition-and-market-context)
3. [System Boundaries and Scope](#system-boundaries-and-scope)
4. [Long-Term Vision and Roadmap](#long-term-vision-and-roadmap)
5. [Version 1 Scope and Deliverables](#version-1-scope-and-deliverables)
6. [Technical Architecture](#technical-architecture)
7. [SDK Architecture and Design](#sdk-architecture-and-design)
8. [Event Architecture](#event-architecture)
9. [Risk Engine Strategy](#risk-engine-strategy)
10. [Policy Engine Architecture](#policy-engine-architecture)
11. [Data Architecture](#data-architecture)
12. [Database Models and Schema](#database-models-and-schema)
13. [Storage Strategy](#storage-strategy)
14. [Scalability and Performance](#scalability-and-performance)
15. [Multi-Tenancy Strategy](#multi-tenancy-strategy)
16. [Security Model](#security-model)
17. [Observability and Monitoring](#observability-and-monitoring)
18. [API Design](#api-design)
19. [Frontend Dashboard](#frontend-dashboard)
20. [Testing Strategy](#testing-strategy)
21. [Data Science Integration](#data-science-integration)
22. [Failure Scenarios and Resilience](#failure-scenarios-and-resilience)
23. [Technical Debt Management](#technical-debt-management)
24. [Engineering Principles](#engineering-principles)
25. [Business Defensibility](#business-defensibility)
26. [Execution Constraints](#execution-constraints)
27. [Technology Stack Decisions](#technology-stack-decisions)

---

## Executive Summary

### Purpose

This document serves as the comprehensive technical blueprint for the AI Control Layer, a middleware and governance platform designed to provide enterprises with centralized oversight, risk management, and policy enforcement for autonomous AI agents.

### Strategic Positioning

The AI Control Layer addresses a critical gap in the AI infrastructure market: the lack of centralized governance and observability for autonomous agents. While numerous tools exist for building and deploying AI agents, few provide enterprise-grade control mechanisms to ensure safe, compliant, and auditable operations.

### Core Value Proposition

- **Centralized Observability**: Monitor all AI agent actions across an organization from a single control plane
- **Risk Management**: Real-time evaluation and scoring of agent outputs to prevent harmful actions
- **Policy Enforcement**: RBAC and rule-based constraints to ensure agents operate within defined boundaries
- **Audit Trail**: Immutable logs of all agent activities for compliance and forensic analysis
- **Developer-Friendly**: SDK-first approach that integrates seamlessly with existing agent frameworks

### Target Market

- AI-native startups deploying autonomous agents in production
- Enterprise development teams building internal AI automation
- Compliance and security teams requiring oversight of AI operations
- Organizations in regulated industries (finance, healthcare, legal) using AI agents

### Competitive Differentiation

Unlike monitoring tools that provide passive observability, the AI Control Layer offers active governance with the ability to intercept, evaluate, and block agent actions in real-time. The SDK-first approach enables early adoption while laying the foundation for a future proxy-based control plane.

---

## Problem Definition and Market Context

### The Rise of Autonomous AI Agents

Organizations are rapidly deploying AI agents to automate complex workflows:

- **Customer Service**: Agents handling support tickets, refunds, and escalations
- **Sales Operations**: Automated lead qualification, email campaigns, and follow-ups
- **Financial Operations**: Transaction approvals, invoice processing, and reconciliation
- **Legal and Compliance**: Contract generation, review, and risk assessment
- **Internal Operations**: Data analysis, reporting, and decision support

### Critical Risks and Challenges

#### 1. Hallucinations and Incorrect Outputs

AI models can generate plausible but incorrect information, leading to:
- Incorrect customer communications
- Faulty financial decisions
- Compliance violations
- Reputational damage

#### 2. Unauthorized Actions

Agents may exceed their intended scope:
- Accessing sensitive data without authorization
- Performing actions beyond their role
- Escalating privileges inappropriately
- Interacting with systems they shouldn't access

#### 3. Cost Overruns

Uncontrolled agent operations can lead to:
- Excessive API calls to LLM providers
- Runaway loops consuming resources
- Inefficient tool usage
- Budget exhaustion

#### 4. Policy Violations

Agents may violate organizational policies:
- Approving refunds beyond authorized limits
- Sharing confidential information
- Making commitments without approval
- Bypassing required workflows

#### 5. Compliance and Regulatory Exposure

Lack of oversight creates compliance risks:
- No audit trail for regulatory review
- Inability to demonstrate due diligence
- Exposure to GDPR, HIPAA, SOC 2 violations
- Lack of explainability for agent decisions

### Current Market Gaps

Existing solutions fall short:

- **LLM Observability Tools**: Focus on model performance, not agent governance
- **APM Solutions**: Monitor infrastructure, not AI-specific risks
- **Agent Frameworks**: Provide building blocks but no control layer
- **Security Tools**: Detect threats but don't understand AI agent context

### Our Solution Approach

The AI Control Layer provides:

1. **Interception Layer**: Wraps all agent actions for visibility
2. **Risk Evaluation**: Real-time scoring of agent outputs
3. **Policy Enforcement**: Active blocking of non-compliant actions
4. **Audit Trail**: Immutable logs for compliance
5. **Centralized Dashboard**: Unified view of all agent operations

---

## System Boundaries and Scope

### In-Scope for V1

#### Core SDK Functionality
- Python SDK that wraps AI agent function calls
- Synchronous event logging to database
- Basic risk scoring using heuristic rules
- Simple policy enforcement (block/allow/warn)
- Local development dashboard for testing

#### Event Management
- Structured event schema for agent actions
- Event types: PROMPT_SENT, RESPONSE_RECEIVED, TOOL_CALLED, TOOL_RESULT, AGENT_LOOP_ITERATION
- Metadata capture: timestamps, agent IDs, token usage, costs
- Synchronous persistence to PostgreSQL

#### Risk Engine (Basic)
- Keyword-based hallucination detection
- Token limit monitoring
- Cost threshold alerts
- Simple confidence scoring
- PII detection using regex patterns

#### Policy Engine (Basic)
- Rule-based policy definitions
- Role-based access control (RBAC) foundations
- Enforcement modes: BLOCK, WARN, ALLOW
- Per-agent policy configuration

#### Dashboard (Development)
- React-based web interface
- Event log viewer
- Risk score visualization
- Policy violation alerts
- Agent activity timeline

### Out-of-Scope for V1

#### Not Building in V1
- Multi-tenant production deployment
- Distributed event streaming (Kafka)
- ML-based risk models
- Proxy-based enforcement
- Enterprise SSO integration
- Advanced analytics and reporting
- Mobile applications
- Third-party integrations (Slack, PagerDuty, etc.)
- Custom agent framework development
- LLM model training or fine-tuning

#### Deferred to Future Versions
- Asynchronous event processing
- Horizontal scaling architecture
- Advanced RBAC with fine-grained permissions
- Compliance report generation
- Multi-region deployment
- High-availability configurations
- Enterprise SLA guarantees

### System Boundaries

#### What We Control
- SDK wrapper logic
- Event schema and storage
- Risk scoring algorithms
- Policy evaluation engine
- Dashboard UI and API

#### What We Don't Control
- Underlying AI models (OpenAI, Anthropic, etc.)
- Agent framework implementations (LangChain, AutoGen)
- Customer's agent business logic
- External tools and APIs agents interact with
- Infrastructure hosting (customer-managed in V1)

### Integration Points

#### Inbound Integrations
- AI agent frameworks (LangChain, AutoGen, custom)
- LLM providers (OpenAI, Anthropic, etc.)
- Customer application code

#### Outbound Integrations
- PostgreSQL database
- Redis for caching (future)
- Logging infrastructure (Loguru)
- Metrics collection (Prometheus)

---

## Long-Term Vision and Roadmap

### 3-5 Year Strategic Vision

Transform from an SDK-first governance tool into a comprehensive AI Control Plane that becomes the standard infrastructure layer for enterprise AI agent deployments.

### Phase 1: SDK-First Foundation (Months 1-6)

**Objective**: Establish core interception and logging capabilities

**Deliverables**:
- Python SDK with agent wrapper functionality
- Event logging to PostgreSQL
- Basic risk scoring heuristics
- Simple policy enforcement
- Development dashboard
- Documentation and examples

**Success Metrics**:
- 10+ early adopter companies
- 100+ agents monitored
- 10,000+ events logged daily
- SDK integration time < 1 hour

**Technical Focus**:
- Correctness over performance
- Clear API design
- Comprehensive testing
- Developer experience

### Phase 2: Centralized Event Collector (Months 7-12)

**Objective**: Scale event processing and add analytics

**Deliverables**:
- Asynchronous event streaming (Redis)
- Centralized event collector service
- Enhanced risk scoring with ML models
- Advanced policy rules engine
- Multi-agent dashboard
- API for programmatic access

**Success Metrics**:
- 50+ companies
- 1,000+ agents monitored
- 1M+ events daily
- <200ms event processing latency

**Technical Focus**:
- Async processing architecture
- Horizontal scalability
- Performance optimization
- Data pipeline reliability

### Phase 3: Proxy-Based Enforcement (Months 13-18)

**Objective**: Enable centralized control without SDK integration

**Deliverables**:
- Proxy layer for LLM API interception
- Network-level agent monitoring
- Centralized policy enforcement
- Multi-tenant architecture
- Enterprise authentication (SSO)
- Compliance reporting

**Success Metrics**:
- 200+ companies
- 10,000+ agents monitored
- 10M+ events daily
- 99.9% uptime SLA

**Technical Focus**:
- Proxy architecture design
- Multi-tenancy isolation
- Enterprise security
- High availability

### Phase 4: Enterprise AI Governance Platform (Months 19-36)

**Objective**: Become the standard AI governance solution

**Deliverables**:
- Advanced analytics and insights
- Anomaly detection with ML
- Automated compliance reports
- Integration marketplace
- Mobile applications
- White-label options
- Global deployment

**Success Metrics**:
- 1,000+ enterprise customers
- 100,000+ agents monitored
- 1B+ events daily
- SOC 2 Type II certified

**Technical Focus**:
- Global infrastructure
- Advanced ML capabilities
- Enterprise integrations
- Platform ecosystem

### Long-Term Competitive Moats

#### Data Network Effects
- Proprietary dataset of agent behaviors
- ML models trained on cross-customer patterns
- Benchmark data for risk scoring

#### Integration Depth
- Deep integrations with all major agent frameworks
- Custom adapters for enterprise systems
- Switching costs increase over time

#### Compliance Positioning
- Become the standard for AI governance
- Required for regulated industries
- Certification and audit partnerships

---

## Version 1 Scope and Deliverables

### Core Components

#### 1. AI Control SDK

**Purpose**: Python library that wraps AI agent operations

**Key Features**:
- Decorator-based function wrapping
- Automatic event capture
- Synchronous logging
- Risk evaluation hooks
- Policy enforcement checks

**API Surface**:
```python
from ai_control import ControlLayer, Agent

# Initialize control layer
control = ControlLayer(api_key="...", org_id="...")

# Wrap agent
@control.monitor(agent_id="customer-support-agent")
def process_customer_request(request):
    # Agent logic here
    pass

# Manual event logging
control.log_event(
    event_type="TOOL_CALLED",
    agent_id="...",
    payload={...}
)
```

**Integration Patterns**:
- LangChain callback handlers
- AutoGen agent wrappers
- Custom agent decorators
- Direct API calls

#### 2. Event Management System

**Purpose**: Capture, structure, and persist agent events

**Event Types**:
- `PROMPT_SENT`: LLM prompt submitted
- `RESPONSE_RECEIVED`: LLM response received
- `TOOL_CALLED`: External tool invoked
- `TOOL_RESULT`: Tool execution result
- `AGENT_LOOP_ITERATION`: Agent reasoning step
- `POLICY_VIOLATION`: Policy rule triggered
- `RISK_ALERT`: Risk threshold exceeded

**Event Schema**:
```json
{
  "event_id": "uuid",
  "event_type": "PROMPT_SENT",
  "timestamp": "ISO8601",
  "agent_id": "string",
  "org_id": "string",
  "payload": {
    "prompt": "string",
    "model": "string",
    "temperature": 0.7,
    "max_tokens": 1000
  },
  "metadata": {
    "user_id": "string",
    "session_id": "string",
    "cost_estimate": 0.05
  },
  "risk_score": 0.3,
  "policy_results": []
}
```

#### 3. Risk Engine (Heuristic)

**Purpose**: Evaluate agent outputs for potential risks

**Risk Categories**:
- Hallucination likelihood
- Policy violation probability
- Toxicity and content safety
- PII leakage detection
- Cost exposure
- Confidence scoring

**Heuristic Rules**:
- Keyword matching for hallucination indicators
- Regex patterns for PII detection
- Token count monitoring
- Cost threshold checks
- Confidence score analysis from LLM responses

**Risk Scoring**:
- Scale: 0.0 (safe) to 1.0 (critical)
- Thresholds: LOW (0-0.3), MEDIUM (0.3-0.7), HIGH (0.7-1.0)
- Composite scoring across multiple dimensions

#### 4. Policy Engine (Basic)

**Purpose**: Enforce organizational rules on agent actions

**Policy Types**:
- **Access Control**: Which agents can perform which actions
- **Resource Limits**: Token budgets, cost caps, rate limits
- **Content Policies**: Prohibited topics, required disclaimers
- **Approval Workflows**: Actions requiring human approval

**Policy Definition**:
```python
{
  "policy_id": "refund-limit",
  "name": "Maximum Refund Amount",
  "description": "Agents cannot approve refunds over $500",
  "rule": {
    "condition": "tool_name == 'approve_refund' AND amount > 500",
    "action": "BLOCK",
    "message": "Refund exceeds authorized limit"
  },
  "applies_to": ["customer-support-agent"],
  "enabled": true
}
```

**Enforcement Modes**:
- `BLOCK`: Prevent action execution
- `WARN`: Log warning but allow
- `REQUIRE_APPROVAL`: Pause for human review

#### 5. Development Dashboard

**Purpose**: Web interface for monitoring and configuration

**Key Views**:
- **Event Stream**: Real-time log of agent activities
- **Agent Overview**: List of monitored agents with status
- **Risk Dashboard**: Visualization of risk scores over time
- **Policy Manager**: Configure and test policies
- **Analytics**: Basic metrics and charts

**Technology**:
- React + Vite for frontend
- FastAPI backend
- WebSocket for real-time updates
- Chart.js for visualizations

### Deliverable Checklist

- [ ] Python SDK package published to PyPI
- [ ] Database schema and migrations
- [ ] FastAPI backend with REST API
- [ ] React dashboard application
- [ ] Documentation website
- [ ] Example integrations (LangChain, AutoGen)
- [ ] Unit and integration tests
- [ ] Docker Compose setup for local development
- [ ] Getting started guide
- [ ] API reference documentation

### Success Criteria for V1

**Functional Requirements**:
- SDK successfully wraps agent calls
- Events logged to database within 200ms
- Risk scores calculated for all events
- Policies enforced before action execution
- Dashboard displays real-time data

**Non-Functional Requirements**:
- SDK adds <50ms latency to agent calls
- Handles 100 events/minute without degradation
- 99% uptime for local deployment
- Clear error messages and logging
- Integration time <1 hour for developers

**User Experience**:
- Simple SDK installation (pip install)
- Minimal configuration required
- Intuitive dashboard interface
- Comprehensive documentation
- Responsive support for early adopters

---

## Technical Architecture

### System Architecture Overview

The AI Control Layer follows a layered architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     AI Agent Application                     │
│              (Customer Code + Agent Framework)               │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                   AI Control SDK (Python)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Wrapper    │  │    Event     │  │    Policy    │     │
│  │   Layer      │  │   Emitter    │  │   Enforcer   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Control Plane API                         │
│                      (FastAPI)                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │    Event     │  │     Risk     │  │    Policy    │     │
│  │   Ingestion  │  │    Engine    │  │    Engine    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Data Layer                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  PostgreSQL  │  │    Redis     │  │   Metrics    │     │
│  │   (Events)   │  │   (Cache)    │  │ (Prometheus) │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                  Dashboard (React + Vite)                    │
└─────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

#### SDK Layer
- Intercept agent function calls
- Capture event data
- Enforce policies synchronously
- Handle errors gracefully
- Minimize performance impact

#### Control Plane API
- Receive events from SDK
- Persist events to database
- Calculate risk scores
- Evaluate policies
- Serve dashboard data

#### Data Layer
- Store events immutably
- Cache frequently accessed data
- Collect metrics
- Support analytics queries

#### Dashboard Layer
- Visualize agent activity
- Configure policies
- Monitor risk scores
- Manage agents and users

### Data Flow

#### Event Capture Flow
1. Agent calls wrapped function
2. SDK intercepts call
3. SDK creates event object
4. SDK sends event to API (sync)
5. API validates event
6. API calculates risk score
7. API evaluates policies
8. API persists to database
9. API returns result to SDK
10. SDK allows/blocks agent action

#### Dashboard Query Flow
1. User opens dashboard
2. Dashboard requests data from API
3. API queries database
4. API aggregates and formats data
5. API returns JSON response
6. Dashboard renders visualizations

### Technology Choices

#### Backend: FastAPI
- Modern async Python framework
- Automatic API documentation (OpenAPI)
- High performance with async/await
- Type hints and validation with Pydantic
- Easy testing and development

#### Database: PostgreSQL
- ACID compliance for audit trails
- JSON support for flexible event payloads
- Mature ecosystem and tooling
- Excellent performance for OLTP workloads
- Support for analytics queries

#### Frontend: React + Vite
- Component-based architecture
- Fast development with hot reload
- Large ecosystem of libraries
- Easy deployment and hosting
- Modern JavaScript tooling

#### Caching: Redis (Future)
- In-memory performance
- Pub/sub for real-time updates
- Simple key-value operations
- Widely adopted and reliable

---


## SDK Architecture and Design

### Design Principles

#### 1. Minimal Intrusion
The SDK should add minimal overhead to agent operations:
- <50ms latency per wrapped call
- No blocking operations in critical path
- Graceful degradation on failures
- Optional async mode for non-critical logging

#### 2. Framework Agnostic
Support multiple agent frameworks without tight coupling:
- Generic wrapper interface
- Framework-specific adapters
- Plugin architecture for extensions
- Standard event schema across frameworks

#### 3. Developer Experience
Make integration as simple as possible:
- Single import statement
- Decorator-based API
- Sensible defaults
- Clear error messages
- Comprehensive examples

#### 4. Fail-Safe Operation
Never break the agent application:
- Catch all SDK exceptions
- Fallback to no-op on errors
- Circuit breaker pattern
- Detailed error logging

### SDK Components

#### Core Wrapper

```python
class ControlLayer:
    """Main SDK interface for agent monitoring and control"""
    
    def __init__(
        self,
        api_key: str,
        org_id: str,
        api_url: str = "http://localhost:8000",
        timeout: float = 5.0,
        fail_open: bool = True
    ):
        self.api_key = api_key
        self.org_id = org_id
        self.api_url = api_url
        self.timeout = timeout
        self.fail_open = fail_open
        self.client = HTTPClient(api_url, api_key)
        
    def monitor(self, agent_id: str, **kwargs):
        """Decorator to wrap agent functions"""
        def decorator(func):
            @wraps(func)
            def wrapper(*args, **kwargs):
                # Pre-execution: log start event
                event = self._create_event(
                    event_type="FUNCTION_START",
                    agent_id=agent_id,
                    payload={"args": args, "kwargs": kwargs}
                )
                
                try:
                    # Check policies before execution
                    policy_result = self._check_policies(event)
                    if policy_result.action == "BLOCK":
                        raise PolicyViolationError(policy_result.message)
                    
                    # Execute function
                    result = func(*args, **kwargs)
                    
                    # Post-execution: log completion event
                    self._log_event(
                        event_type="FUNCTION_COMPLETE",
                        agent_id=agent_id,
                        payload={"result": result}
                    )
                    
                    return result
                    
                except Exception as e:
                    # Log error event
                    self._log_event(
                        event_type="FUNCTION_ERROR",
                        agent_id=agent_id,
                        payload={"error": str(e)}
                    )
                    raise
                    
            return wrapper
        return decorator
```

#### Event Emitter

```python
class EventEmitter:
    """Handles event creation and transmission"""
    
    def emit(self, event: Event) -> bool:
        """Send event to control plane"""
        try:
            response = self.client.post(
                "/api/v1/events",
                json=event.dict(),
                timeout=self.timeout
            )
            return response.status_code == 200
        except Exception as e:
            logger.error(f"Failed to emit event: {e}")
            if not self.fail_open:
                raise
            return False
```

#### Policy Enforcer

```python
class PolicyEnforcer:
    """Evaluates policies before agent actions"""
    
    def check(self, event: Event) -> PolicyResult:
        """Check if event violates any policies"""
        try:
            response = self.client.post(
                "/api/v1/policies/evaluate",
                json=event.dict(),
                timeout=self.timeout
            )
            return PolicyResult(**response.json())
        except Exception as e:
            logger.error(f"Policy check failed: {e}")
            if self.fail_open:
                return PolicyResult(action="ALLOW")
            raise
```

### Framework Integrations

#### LangChain Integration

```python
from langchain.callbacks.base import BaseCallbackHandler

class ControlLayerCallback(BaseCallbackHandler):
    """LangChain callback for AI Control Layer"""
    
    def __init__(self, control: ControlLayer, agent_id: str):
        self.control = control
        self.agent_id = agent_id
    
    def on_llm_start(self, serialized, prompts, **kwargs):
        """Called when LLM starts"""
        self.control.log_event(
            event_type="PROMPT_SENT",
            agent_id=self.agent_id,
            payload={"prompts": prompts, "model": serialized}
        )
    
    def on_llm_end(self, response, **kwargs):
        """Called when LLM completes"""
        self.control.log_event(
            event_type="RESPONSE_RECEIVED",
            agent_id=self.agent_id,
            payload={"response": response.dict()}
        )
    
    def on_tool_start(self, tool, input_str, **kwargs):
        """Called when tool starts"""
        self.control.log_event(
            event_type="TOOL_CALLED",
            agent_id=self.agent_id,
            payload={"tool": tool.name, "input": input_str}
        )
```

#### AutoGen Integration

```python
from autogen import Agent

class ControlledAgent(Agent):
    """AutoGen agent with AI Control Layer"""
    
    def __init__(self, control: ControlLayer, agent_id: str, **kwargs):
        super().__init__(**kwargs)
        self.control = control
        self.agent_id = agent_id
    
    def send(self, message, recipient, **kwargs):
        """Override send to log events"""
        self.control.log_event(
            event_type="AGENT_MESSAGE",
            agent_id=self.agent_id,
            payload={"message": message, "recipient": recipient.name}
        )
        return super().send(message, recipient, **kwargs)
```

### SDK Configuration

#### Environment Variables

```bash
AI_CONTROL_API_KEY=your_api_key_here
AI_CONTROL_ORG_ID=your_org_id
AI_CONTROL_API_URL=http://localhost:8000
AI_CONTROL_TIMEOUT=5.0
AI_CONTROL_FAIL_OPEN=true
AI_CONTROL_LOG_LEVEL=INFO
```

#### Configuration File

```yaml
# ai_control_config.yaml
api_key: ${AI_CONTROL_API_KEY}
org_id: ${AI_CONTROL_ORG_ID}
api_url: http://localhost:8000
timeout: 5.0
fail_open: true

# Event filtering
events:
  include:
    - PROMPT_SENT
    - RESPONSE_RECEIVED
    - TOOL_CALLED
  exclude:
    - DEBUG_*

# Sampling (for high-volume scenarios)
sampling:
  enabled: false
  rate: 0.1  # Log 10% of events

# Batching (future)
batching:
  enabled: false
  max_size: 100
  max_wait_ms: 1000
```

### Error Handling

#### Circuit Breaker Pattern

```python
class CircuitBreaker:
    """Prevent cascading failures"""
    
    def __init__(self, failure_threshold: int = 5, timeout: int = 60):
        self.failure_threshold = failure_threshold
        self.timeout = timeout
        self.failures = 0
        self.last_failure_time = None
        self.state = "CLOSED"  # CLOSED, OPEN, HALF_OPEN
    
    def call(self, func, *args, **kwargs):
        if self.state == "OPEN":
            if time.time() - self.last_failure_time > self.timeout:
                self.state = "HALF_OPEN"
            else:
                raise CircuitBreakerOpenError()
        
        try:
            result = func(*args, **kwargs)
            if self.state == "HALF_OPEN":
                self.state = "CLOSED"
                self.failures = 0
            return result
        except Exception as e:
            self.failures += 1
            self.last_failure_time = time.time()
            if self.failures >= self.failure_threshold:
                self.state = "OPEN"
            raise
```

### Performance Optimization

#### Async Event Emission (Future)

```python
import asyncio
from queue import Queue

class AsyncEventEmitter:
    """Non-blocking event emission"""
    
    def __init__(self):
        self.queue = Queue()
        self.worker_thread = Thread(target=self._worker)
        self.worker_thread.start()
    
    def emit_async(self, event: Event):
        """Add event to queue for async processing"""
        self.queue.put(event)
    
    def _worker(self):
        """Background worker to send events"""
        while True:
            event = self.queue.get()
            try:
                self._send_event(event)
            except Exception as e:
                logger.error(f"Failed to send event: {e}")
            finally:
                self.queue.task_done()
```

---

## Event Architecture

### Event Schema Design

#### Core Event Structure

Every event follows a canonical schema:

```python
from pydantic import BaseModel, Field
from typing import Dict, Any, Optional, List
from datetime import datetime
from enum import Enum

class EventType(str, Enum):
    PROMPT_SENT = "PROMPT_SENT"
    RESPONSE_RECEIVED = "RESPONSE_RECEIVED"
    TOOL_CALLED = "TOOL_CALLED"
    TOOL_RESULT = "TOOL_RESULT"
    AGENT_LOOP_ITERATION = "AGENT_LOOP_ITERATION"
    POLICY_VIOLATION = "POLICY_VIOLATION"
    RISK_ALERT = "RISK_ALERT"
    FUNCTION_START = "FUNCTION_START"
    FUNCTION_COMPLETE = "FUNCTION_COMPLETE"
    FUNCTION_ERROR = "FUNCTION_ERROR"

class Event(BaseModel):
    # Required fields
    event_id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    event_type: EventType
    timestamp: datetime = Field(default_factory=datetime.utcnow)
    agent_id: str
    org_id: str
    
    # Event-specific data
    payload: Dict[str, Any]
    
    # Metadata
    metadata: Optional[Dict[str, Any]] = Field(default_factory=dict)
    
    # Computed fields (added by control plane)
    risk_score: Optional[float] = None
    policy_results: Optional[List[Dict]] = None
    
    # Tracing
    trace_id: Optional[str] = None
    parent_event_id: Optional[str] = None
    
    # Cost tracking
    estimated_cost: Optional[float] = None
    token_count: Optional[int] = None
    
    class Config:
        json_encoders = {
            datetime: lambda v: v.isoformat()
        }
```

### Event Types and Payloads

#### PROMPT_SENT

Captured when an LLM prompt is submitted:

```json
{
  "event_type": "PROMPT_SENT",
  "agent_id": "customer-support-agent-001",
  "payload": {
    "prompt": "You are a helpful customer support agent...",
    "model": "gpt-4",
    "temperature": 0.7,
    "max_tokens": 1000,
    "system_message": "...",
    "user_message": "..."
  },
  "metadata": {
    "user_id": "user_123",
    "session_id": "session_456",
    "request_id": "req_789"
  },
  "estimated_cost": 0.03,
  "token_count": 500
}
```

#### RESPONSE_RECEIVED

Captured when LLM returns a response:

```json
{
  "event_type": "RESPONSE_RECEIVED",
  "agent_id": "customer-support-agent-001",
  "payload": {
    "response": "I'd be happy to help you with...",
    "model": "gpt-4",
    "finish_reason": "stop",
    "usage": {
      "prompt_tokens": 500,
      "completion_tokens": 150,
      "total_tokens": 650
    }
  },
  "parent_event_id": "event_prompt_sent_id",
  "estimated_cost": 0.05,
  "token_count": 650,
  "risk_score": 0.2
}
```

#### TOOL_CALLED

Captured when agent invokes an external tool:

```json
{
  "event_type": "TOOL_CALLED",
  "agent_id": "customer-support-agent-001",
  "payload": {
    "tool_name": "approve_refund",
    "tool_description": "Approve a customer refund",
    "arguments": {
      "customer_id": "cust_123",
      "amount": 49.99,
      "reason": "Product defect"
    }
  },
  "metadata": {
    "requires_approval": false,
    "risk_level": "low"
  }
}
```

#### TOOL_RESULT

Captured when tool execution completes:

```json
{
  "event_type": "TOOL_RESULT",
  "agent_id": "customer-support-agent-001",
  "payload": {
    "tool_name": "approve_refund",
    "result": {
      "success": true,
      "refund_id": "ref_789",
      "message": "Refund approved"
    },
    "execution_time_ms": 234
  },
  "parent_event_id": "event_tool_called_id"
}
```

#### AGENT_LOOP_ITERATION

Captured for each reasoning step in agent loop:

```json
{
  "event_type": "AGENT_LOOP_ITERATION",
  "agent_id": "customer-support-agent-001",
  "payload": {
    "iteration": 3,
    "thought": "I need to check the customer's order history",
    "action": "query_database",
    "action_input": {"customer_id": "cust_123"},
    "observation": "Customer has 5 previous orders"
  },
  "metadata": {
    "loop_type": "react",
    "max_iterations": 10
  }
}
```

#### POLICY_VIOLATION

Captured when a policy is violated:

```json
{
  "event_type": "POLICY_VIOLATION",
  "agent_id": "customer-support-agent-001",
  "payload": {
    "policy_id": "refund-limit-policy",
    "policy_name": "Maximum Refund Amount",
    "violation_type": "THRESHOLD_EXCEEDED",
    "details": {
      "limit": 500,
      "attempted": 750,
      "action": "BLOCKED"
    }
  },
  "metadata": {
    "severity": "HIGH",
    "requires_review": true
  }
}
```

### Event Lifecycle

#### 1. Creation
- Event created in SDK when agent action occurs
- Assigned unique event_id
- Timestamp captured
- Payload populated with action data

#### 2. Enrichment
- Metadata added (user context, session info)
- Trace ID for request correlation
- Parent event ID for causality tracking

#### 3. Transmission
- Event serialized to JSON
- Sent to control plane API
- Retry logic for transient failures

#### 4. Processing
- Event validated against schema
- Risk score calculated
- Policies evaluated
- Results added to event

#### 5. Persistence
- Event stored in PostgreSQL
- Indexed for fast queries
- Immutable once written

#### 6. Analytics
- Event available for dashboard queries
- Aggregated for metrics
- Used for ML model training (future)

### Event Storage Strategy

#### Database Schema

```sql
CREATE TABLE events (
    event_id UUID PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    agent_id VARCHAR(255) NOT NULL,
    org_id VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB,
    risk_score FLOAT,
    policy_results JSONB,
    trace_id VARCHAR(255),
    parent_event_id UUID,
    estimated_cost FLOAT,
    token_count INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_events_agent_id ON events(agent_id);
CREATE INDEX idx_events_org_id ON events(org_id);
CREATE INDEX idx_events_timestamp ON events(timestamp DESC);
CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_trace_id ON events(trace_id);
CREATE INDEX idx_events_risk_score ON events(risk_score) WHERE risk_score > 0.7;

-- Composite indexes for dashboard queries
CREATE INDEX idx_events_org_agent_time ON events(org_id, agent_id, timestamp DESC);
```

#### Partitioning Strategy (Future)

```sql
-- Partition by month for time-series data
CREATE TABLE events_2024_01 PARTITION OF events
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE events_2024_02 PARTITION OF events
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

### Event Query Patterns

#### Recent Events for Agent

```sql
SELECT * FROM events
WHERE agent_id = 'customer-support-agent-001'
  AND timestamp > NOW() - INTERVAL '1 hour'
ORDER BY timestamp DESC
LIMIT 100;
```

#### High-Risk Events

```sql
SELECT * FROM events
WHERE org_id = 'org_123'
  AND risk_score > 0.7
  AND timestamp > NOW() - INTERVAL '24 hours'
ORDER BY risk_score DESC;
```

#### Policy Violations

```sql
SELECT * FROM events
WHERE event_type = 'POLICY_VIOLATION'
  AND org_id = 'org_123'
  AND timestamp > NOW() - INTERVAL '7 days';
```

#### Cost Analysis

```sql
SELECT 
    agent_id,
    DATE(timestamp) as date,
    SUM(estimated_cost) as total_cost,
    SUM(token_count) as total_tokens,
    COUNT(*) as event_count
FROM events
WHERE org_id = 'org_123'
  AND timestamp > NOW() - INTERVAL '30 days'
GROUP BY agent_id, DATE(timestamp)
ORDER BY total_cost DESC;
```

### Event Streaming (Future)

#### Redis Pub/Sub

```python
import redis

class EventStreamer:
    def __init__(self, redis_url: str):
        self.redis = redis.from_url(redis_url)
    
    def publish(self, event: Event):
        """Publish event to Redis channel"""
        channel = f"events:{event.org_id}:{event.agent_id}"
        self.redis.publish(channel, event.json())
    
    def subscribe(self, org_id: str, agent_id: str):
        """Subscribe to event stream"""
        channel = f"events:{org_id}:{agent_id}"
        pubsub = self.redis.pubsub()
        pubsub.subscribe(channel)
        return pubsub
```

#### Kafka Integration (Future)

```python
from aiokafka import AIOKafkaProducer

class KafkaEventStreamer:
    def __init__(self, bootstrap_servers: str):
        self.producer = AIOKafkaProducer(
            bootstrap_servers=bootstrap_servers
        )
    
    async def publish(self, event: Event):
        """Publish event to Kafka topic"""
        topic = f"ai-control-events-{event.org_id}"
        await self.producer.send(
            topic,
            key=event.agent_id.encode(),
            value=event.json().encode()
        )
```

---


## Risk Engine Strategy

### Overview

The Risk Engine evaluates agent outputs to identify potential issues before they cause harm. In V1, we use heuristic-based scoring; future versions will incorporate ML models trained on historical data.

### Risk Dimensions

#### 1. Hallucination Detection

**Indicators**:
- Hedging language ("I think", "maybe", "possibly")
- Contradictory statements within response
- Lack of specific details when expected
- Confidence scores below threshold
- Unusual response patterns

**Heuristic Rules**:
```python
def detect_hallucination(response: str, confidence: float) -> float:
    """Calculate hallucination risk score"""
    score = 0.0
    
    # Check for hedging language
    hedging_words = ["i think", "maybe", "possibly", "might", "perhaps"]
    hedge_count = sum(1 for word in hedging_words if word in response.lower())
    score += min(hedge_count * 0.1, 0.3)
    
    # Check confidence score
    if confidence < 0.7:
        score += (0.7 - confidence) * 0.5
    
    # Check response length (very short or very long can indicate issues)
    word_count = len(response.split())
    if word_count < 10:
        score += 0.2
    elif word_count > 500:
        score += 0.1
    
    return min(score, 1.0)
```

#### 2. Policy Violation Probability

**Indicators**:
- Actions exceeding authorized limits
- Access to restricted resources
- Prohibited content or topics
- Missing required approvals

**Heuristic Rules**:
```python
def check_policy_violation(event: Event, policies: List[Policy]) -> float:
    """Calculate policy violation risk"""
    violations = []
    
    for policy in policies:
        if policy.applies_to_agent(event.agent_id):
            result = policy.evaluate(event)
            if result.violated:
                violations.append(result)
    
    if not violations:
        return 0.0
    
    # Weight by severity
    severity_weights = {"LOW": 0.3, "MEDIUM": 0.6, "HIGH": 0.9, "CRITICAL": 1.0}
    max_severity = max(severity_weights[v.severity] for v in violations)
    
    return max_severity
```

#### 3. Toxicity and Content Safety

**Indicators**:
- Offensive language
- Discriminatory content
- Threats or violence
- Inappropriate topics

**Heuristic Rules**:
```python
def detect_toxicity(text: str) -> float:
    """Calculate toxicity score"""
    score = 0.0
    
    # Offensive words list (simplified)
    offensive_words = load_offensive_words_list()
    offensive_count = sum(1 for word in offensive_words if word in text.lower())
    score += min(offensive_count * 0.2, 0.6)
    
    # ALL CAPS (aggressive tone)
    caps_ratio = sum(1 for c in text if c.isupper()) / max(len(text), 1)
    if caps_ratio > 0.5:
        score += 0.2
    
    # Excessive punctuation (!!!, ???)
    exclamation_count = text.count('!')
    if exclamation_count > 3:
        score += 0.1
    
    return min(score, 1.0)
```

#### 4. PII Leakage Detection

**Indicators**:
- Email addresses
- Phone numbers
- Social Security Numbers
- Credit card numbers
- Physical addresses

**Heuristic Rules**:
```python
import re

def detect_pii(text: str) -> float:
    """Calculate PII leakage risk"""
    pii_patterns = {
        'email': r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b',
        'phone': r'\b\d{3}[-.]?\d{3}[-.]?\d{4}\b',
        'ssn': r'\b\d{3}-\d{2}-\d{4}\b',
        'credit_card': r'\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b'
    }
    
    pii_found = []
    for pii_type, pattern in pii_patterns.items():
        matches = re.findall(pattern, text)
        if matches:
            pii_found.append((pii_type, len(matches)))
    
    if not pii_found:
        return 0.0
    
    # High risk if any PII detected
    return 0.8 + (len(pii_found) * 0.05)
```

#### 5. Cost Exposure

**Indicators**:
- High token usage
- Expensive model calls
- Runaway loops
- Excessive API calls

**Heuristic Rules**:
```python
def calculate_cost_risk(event: Event, budget: float) -> float:
    """Calculate cost exposure risk"""
    if event.estimated_cost is None:
        return 0.0
    
    # Check against budget
    cost_ratio = event.estimated_cost / budget
    
    if cost_ratio < 0.5:
        return 0.0
    elif cost_ratio < 0.8:
        return 0.3
    elif cost_ratio < 1.0:
        return 0.6
    else:
        return 1.0  # Over budget
```

### Composite Risk Scoring

Combine multiple risk dimensions into a single score:

```python
class RiskEngine:
    """Calculate composite risk scores for agent events"""
    
    def __init__(self, weights: Dict[str, float] = None):
        self.weights = weights or {
            'hallucination': 0.3,
            'policy_violation': 0.3,
            'toxicity': 0.2,
            'pii_leakage': 0.15,
            'cost_exposure': 0.05
        }
    
    def calculate_risk(self, event: Event, context: Dict) -> RiskScore:
        """Calculate composite risk score"""
        scores = {
            'hallucination': detect_hallucination(
                event.payload.get('response', ''),
                event.payload.get('confidence', 1.0)
            ),
            'policy_violation': check_policy_violation(
                event,
                context.get('policies', [])
            ),
            'toxicity': detect_toxicity(
                event.payload.get('response', '')
            ),
            'pii_leakage': detect_pii(
                event.payload.get('response', '')
            ),
            'cost_exposure': calculate_cost_risk(
                event,
                context.get('budget', float('inf'))
            )
        }
        
        # Weighted average
        composite_score = sum(
            scores[dim] * self.weights[dim]
            for dim in scores
        )
        
        return RiskScore(
            composite=composite_score,
            dimensions=scores,
            severity=self._get_severity(composite_score),
            recommendations=self._get_recommendations(scores)
        )
    
    def _get_severity(self, score: float) -> str:
        """Map score to severity level"""
        if score < 0.3:
            return "LOW"
        elif score < 0.7:
            return "MEDIUM"
        else:
            return "HIGH"
    
    def _get_recommendations(self, scores: Dict[str, float]) -> List[str]:
        """Generate actionable recommendations"""
        recommendations = []
        
        if scores['hallucination'] > 0.5:
            recommendations.append("Consider requesting clarification or additional context")
        
        if scores['policy_violation'] > 0.7:
            recommendations.append("Block action and require human approval")
        
        if scores['toxicity'] > 0.6:
            recommendations.append("Review content for appropriateness")
        
        if scores['pii_leakage'] > 0.5:
            recommendations.append("Redact PII before logging or transmission")
        
        if scores['cost_exposure'] > 0.8:
            recommendations.append("Approaching budget limit - consider throttling")
        
        return recommendations
```

### Risk Score Storage

```python
class RiskScore(BaseModel):
    """Risk score result"""
    composite: float = Field(ge=0.0, le=1.0)
    dimensions: Dict[str, float]
    severity: str  # LOW, MEDIUM, HIGH
    recommendations: List[str]
    calculated_at: datetime = Field(default_factory=datetime.utcnow)
```

### Future ML-Based Risk Scoring

#### Training Data Collection

```python
# Collect labeled examples
training_data = {
    'event_features': [
        # Extract features from events
        {
            'response_length': 150,
            'confidence': 0.85,
            'token_count': 200,
            'hedge_word_count': 2,
            'caps_ratio': 0.05,
            # ... more features
        }
    ],
    'labels': [
        # Human-labeled risk scores
        {'hallucination_risk': 0.3, 'overall_risk': 0.4}
    ]
}
```

#### Model Training

```python
from sklearn.ensemble import RandomForestRegressor

class MLRiskEngine:
    """ML-based risk scoring (future)"""
    
    def __init__(self):
        self.model = RandomForestRegressor()
        self.feature_extractor = FeatureExtractor()
    
    def train(self, events: List[Event], labels: List[float]):
        """Train risk model on historical data"""
        features = [self.feature_extractor.extract(e) for e in events]
        self.model.fit(features, labels)
    
    def predict(self, event: Event) -> float:
        """Predict risk score for new event"""
        features = self.feature_extractor.extract(event)
        return self.model.predict([features])[0]
```

---

## Policy Engine Architecture

### Overview

The Policy Engine enforces organizational rules and constraints on agent actions. It evaluates events against defined policies and determines whether to allow, block, or require approval for actions.

### Policy Types

#### 1. Access Control Policies

Define which agents can perform which actions:

```python
class AccessControlPolicy(BaseModel):
    policy_id: str
    name: str
    description: str
    agent_ids: List[str]  # Agents this policy applies to
    allowed_actions: List[str]
    denied_actions: List[str]
    
    def evaluate(self, event: Event) -> PolicyResult:
        """Check if agent is authorized for action"""
        if event.agent_id not in self.agent_ids:
            return PolicyResult(action="ALLOW", reason="Policy does not apply")
        
        action = event.payload.get('action')
        
        if action in self.denied_actions:
            return PolicyResult(
                action="BLOCK",
                reason=f"Action '{action}' is explicitly denied"
            )
        
        if self.allowed_actions and action not in self.allowed_actions:
            return PolicyResult(
                action="BLOCK",
                reason=f"Action '{action}' is not in allowed list"
            )
        
        return PolicyResult(action="ALLOW")
```

#### 2. Resource Limit Policies

Enforce limits on resource usage:

```python
class ResourceLimitPolicy(BaseModel):
    policy_id: str
    name: str
    limits: Dict[str, Any]  # e.g., {'max_tokens': 10000, 'max_cost': 100}
    time_window: str  # e.g., 'hourly', 'daily', 'monthly'
    
    def evaluate(self, event: Event, usage_tracker: UsageTracker) -> PolicyResult:
        """Check if action would exceed resource limits"""
        current_usage = usage_tracker.get_usage(
            agent_id=event.agent_id,
            time_window=self.time_window
        )
        
        # Check token limit
        if 'max_tokens' in self.limits:
            projected_tokens = current_usage['tokens'] + event.token_count
            if projected_tokens > self.limits['max_tokens']:
                return PolicyResult(
                    action="BLOCK",
                    reason=f"Would exceed token limit: {projected_tokens}/{self.limits['max_tokens']}"
                )
        
        # Check cost limit
        if 'max_cost' in self.limits:
            projected_cost = current_usage['cost'] + event.estimated_cost
            if projected_cost > self.limits['max_cost']:
                return PolicyResult(
                    action="BLOCK",
                    reason=f"Would exceed cost limit: ${projected_cost:.2f}/${self.limits['max_cost']}"
                )
        
        return PolicyResult(action="ALLOW")
```

#### 3. Content Policies

Enforce rules about content and topics:

```python
class ContentPolicy(BaseModel):
    policy_id: str
    name: str
    prohibited_topics: List[str]
    required_disclaimers: List[str]
    max_response_length: Optional[int] = None
    
    def evaluate(self, event: Event) -> PolicyResult:
        """Check content compliance"""
        response = event.payload.get('response', '')
        
        # Check prohibited topics
        for topic in self.prohibited_topics:
            if topic.lower() in response.lower():
                return PolicyResult(
                    action="BLOCK",
                    reason=f"Response contains prohibited topic: {topic}"
                )
        
        # Check required disclaimers
        for disclaimer in self.required_disclaimers:
            if disclaimer not in response:
                return PolicyResult(
                    action="WARN",
                    reason=f"Response missing required disclaimer: {disclaimer}"
                )
        
        # Check length
        if self.max_response_length:
            if len(response) > self.max_response_length:
                return PolicyResult(
                    action="WARN",
                    reason=f"Response exceeds maximum length: {len(response)}/{self.max_response_length}"
                )
        
        return PolicyResult(action="ALLOW")
```

#### 4. Approval Workflow Policies

Require human approval for certain actions:

```python
class ApprovalPolicy(BaseModel):
    policy_id: str
    name: str
    requires_approval_if: Dict[str, Any]  # Conditions requiring approval
    approvers: List[str]  # User IDs who can approve
    timeout_minutes: int = 60
    
    def evaluate(self, event: Event) -> PolicyResult:
        """Check if action requires approval"""
        # Check conditions
        for condition, threshold in self.requires_approval_if.items():
            if condition == 'amount' and 'amount' in event.payload:
                if event.payload['amount'] > threshold:
                    return PolicyResult(
                        action="REQUIRE_APPROVAL",
                        reason=f"Amount ${event.payload['amount']} exceeds threshold ${threshold}",
                        metadata={
                            'approvers': self.approvers,
                            'timeout_minutes': self.timeout_minutes
                        }
                    )
            
            if condition == 'risk_score' and event.risk_score:
                if event.risk_score > threshold:
                    return PolicyResult(
                        action="REQUIRE_APPROVAL",
                        reason=f"Risk score {event.risk_score} exceeds threshold {threshold}",
                        metadata={
                            'approvers': self.approvers,
                            'timeout_minutes': self.timeout_minutes
                        }
                    )
        
        return PolicyResult(action="ALLOW")
```

### Policy Evaluation Engine

```python
class PolicyEngine:
    """Evaluate events against all applicable policies"""
    
    def __init__(self, policy_store: PolicyStore):
        self.policy_store = policy_store
    
    def evaluate(self, event: Event, context: Dict = None) -> PolicyEvaluationResult:
        """Evaluate event against all policies"""
        context = context or {}
        
        # Get applicable policies
        policies = self.policy_store.get_policies_for_agent(event.agent_id)
        
        results = []
        for policy in policies:
            result = policy.evaluate(event, context)
            results.append({
                'policy_id': policy.policy_id,
                'policy_name': policy.name,
                'action': result.action,
                'reason': result.reason,
                'metadata': result.metadata
            })
        
        # Determine final action (most restrictive wins)
        final_action = self._determine_final_action(results)
        
        return PolicyEvaluationResult(
            action=final_action,
            policy_results=results,
            evaluated_at=datetime.utcnow()
        )
    
    def _determine_final_action(self, results: List[Dict]) -> str:
        """Determine final action from multiple policy results"""
        actions = [r['action'] for r in results]
        
        # Priority: BLOCK > REQUIRE_APPROVAL > WARN > ALLOW
        if "BLOCK" in actions:
            return "BLOCK"
        elif "REQUIRE_APPROVAL" in actions:
            return "REQUIRE_APPROVAL"
        elif "WARN" in actions:
            return "WARN"
        else:
            return "ALLOW"
```

### Policy Configuration

#### YAML Configuration

```yaml
policies:
  - policy_id: refund-limit
    type: approval
    name: Refund Approval Policy
    description: Require approval for refunds over $500
    applies_to:
      - customer-support-agent
    requires_approval_if:
      amount: 500
    approvers:
      - manager@company.com
      - supervisor@company.com
    timeout_minutes: 60
    enabled: true

  - policy_id: token-budget
    type: resource_limit
    name: Daily Token Budget
    description: Limit token usage to 100k per day
    applies_to:
      - "*"  # All agents
    limits:
      max_tokens: 100000
      max_cost: 50.00
    time_window: daily
    enabled: true

  - policy_id: content-safety
    type: content
    name: Content Safety Policy
    description: Prohibit sensitive topics
    applies_to:
      - "*"
    prohibited_topics:
      - politics
      - religion
      - medical advice
    required_disclaimers:
      - "This is AI-generated content"
    enabled: true
```

### Policy Storage

```sql
CREATE TABLE policies (
    policy_id VARCHAR(255) PRIMARY KEY,
    org_id VARCHAR(255) NOT NULL,
    policy_type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    config JSONB NOT NULL,
    applies_to JSONB NOT NULL,  -- List of agent IDs or patterns
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_policies_org_id ON policies(org_id);
CREATE INDEX idx_policies_enabled ON policies(enabled) WHERE enabled = true;
```

### Policy Audit Trail

```sql
CREATE TABLE policy_evaluations (
    evaluation_id UUID PRIMARY KEY,
    event_id UUID REFERENCES events(event_id),
    policy_id VARCHAR(255) REFERENCES policies(policy_id),
    action VARCHAR(50) NOT NULL,  -- ALLOW, BLOCK, WARN, REQUIRE_APPROVAL
    reason TEXT,
    metadata JSONB,
    evaluated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_policy_eval_event ON policy_evaluations(event_id);
CREATE INDEX idx_policy_eval_policy ON policy_evaluations(policy_id);
```

---

## Data Architecture

### Design Principles

#### 1. Immutability
All events are immutable once written. This ensures:
- Audit trail integrity
- Compliance with regulatory requirements
- Ability to replay and analyze historical data
- Trust in the system

#### 2. Extensibility
Schema design supports future additions:
- JSONB columns for flexible payloads
- Metadata fields for custom attributes
- Versioned event schemas
- Backward compatibility

#### 3. Performance
Optimize for common query patterns:
- Indexes on frequently filtered columns
- Partitioning for time-series data
- Denormalization where appropriate
- Caching for hot data

#### 4. Multi-Tenancy Ready
Design for future multi-tenant deployment:
- org_id in all tables
- Row-level security policies
- Tenant isolation
- Separate schemas or databases per tenant (future)

### Core Data Models

#### Organizations

```sql
CREATE TABLE organizations (
    org_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(50) NOT NULL,  -- free, pro, enterprise
    settings JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Agents

```sql
CREATE TABLE agents (
    agent_id VARCHAR(255) PRIMARY KEY,
    org_id VARCHAR(255) REFERENCES organizations(org_id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    framework VARCHAR(50),  -- langchain, autogen, custom
    config JSONB,
    status VARCHAR(50) DEFAULT 'active',  -- active, paused, archived
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_agents_org_id ON agents(org_id);
CREATE INDEX idx_agents_status ON agents(status);
```

#### Events (Already Defined)

```sql
CREATE TABLE events (
    event_id UUID PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    agent_id VARCHAR(255) REFERENCES agents(agent_id),
    org_id VARCHAR(255) REFERENCES organizations(org_id),
    payload JSONB NOT NULL,
    metadata JSONB,
    risk_score FLOAT,
    policy_results JSONB,
    trace_id VARCHAR(255),
    parent_event_id UUID REFERENCES events(event_id),
    estimated_cost FLOAT,
    token_count INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Risk Scores

```sql
CREATE TABLE risk_scores (
    score_id UUID PRIMARY KEY,
    event_id UUID REFERENCES events(event_id),
    composite_score FLOAT NOT NULL,
    hallucination_score FLOAT,
    policy_violation_score FLOAT,
    toxicity_score FLOAT,
    pii_leakage_score FLOAT,
    cost_exposure_score FLOAT,
    severity VARCHAR(50) NOT NULL,
    recommendations JSONB,
    calculated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_risk_scores_event ON risk_scores(event_id);
CREATE INDEX idx_risk_scores_severity ON risk_scores(severity);
```

#### Policies (Already Defined)

```sql
CREATE TABLE policies (
    policy_id VARCHAR(255) PRIMARY KEY,
    org_id VARCHAR(255) REFERENCES organizations(org_id),
    policy_type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    config JSONB NOT NULL,
    applies_to JSONB NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Users

```sql
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    org_id VARCHAR(255) REFERENCES organizations(org_id),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) NOT NULL,  -- admin, developer, viewer
    permissions JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_users_org_id ON users(org_id);
CREATE INDEX idx_users_email ON users(email);
```

#### API Keys

```sql
CREATE TABLE api_keys (
    key_id VARCHAR(255) PRIMARY KEY,
    org_id VARCHAR(255) REFERENCES organizations(org_id),
    key_hash VARCHAR(255) NOT NULL,  -- Hashed API key
    name VARCHAR(255),
    permissions JSONB,
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_api_keys_org_id ON api_keys(org_id);
CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);
```

### Data Relationships

```
organizations (1) ----< (N) agents
organizations (1) ----< (N) users
organizations (1) ----< (N) policies
organizations (1) ----< (N) api_keys

agents (1) ----< (N) events

events (1) ----< (1) risk_scores
events (1) ----< (N) policy_evaluations

policies (1) ----< (N) policy_evaluations
```

### Data Retention

#### Retention Policies

```python
class RetentionPolicy:
    """Define data retention rules"""
    
    # Events older than 90 days archived to cold storage
    EVENT_HOT_STORAGE_DAYS = 90
    
    # Events older than 1 year deleted (unless compliance requires longer)
    EVENT_MAX_RETENTION_DAYS = 365
    
    # Risk scores retained for 180 days
    RISK_SCORE_RETENTION_DAYS = 180
    
    # Policy evaluations retained for 1 year
    POLICY_EVAL_RETENTION_DAYS = 365
```

#### Archival Process

```sql
-- Move old events to archive table
INSERT INTO events_archive
SELECT * FROM events
WHERE timestamp < NOW() - INTERVAL '90 days';

DELETE FROM events
WHERE timestamp < NOW() - INTERVAL '90 days';
```

---


## Database Models and Schema

### SQLAlchemy Models

#### Organization Model

```python
from sqlalchemy import Column, String, TIMESTAMP, JSON
from sqlalchemy.orm import relationship
from datetime import datetime

class Organization(Base):
    __tablename__ = 'organizations'
    
    org_id = Column(String(255), primary_key=True)
    name = Column(String(255), nullable=False)
    plan = Column(String(50), nullable=False, default='free')
    settings = Column(JSON)
    created_at = Column(TIMESTAMP, default=datetime.utcnow)
    updated_at = Column(TIMESTAMP, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    # Relationships
    agents = relationship("Agent", back_populates="organization")
    users = relationship("User", back_populates="organization")
    policies = relationship("Policy", back_populates="organization")
```

#### Agent Model

```python
class Agent(Base):
    __tablename__ = 'agents'
    
    agent_id = Column(String(255), primary_key=True)
    org_id = Column(String(255), ForeignKey('organizations.org_id'), nullable=False)
    name = Column(String(255), nullable=False)
    description = Column(Text)
    framework = Column(String(50))
    config = Column(JSON)
    status = Column(String(50), default='active')
    created_at = Column(TIMESTAMP, default=datetime.utcnow)
    updated_at = Column(TIMESTAMP, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    # Relationships
    organization = relationship("Organization", back_populates="agents")
    events = relationship("Event", back_populates="agent")
```

#### Event Model

```python
from sqlalchemy.dialects.postgresql import UUID, JSONB
import uuid

class Event(Base):
    __tablename__ = 'events'
    
    event_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    event_type = Column(String(50), nullable=False)
    timestamp = Column(TIMESTAMP, nullable=False, default=datetime.utcnow)
    agent_id = Column(String(255), ForeignKey('agents.agent_id'), nullable=False)
    org_id = Column(String(255), ForeignKey('organizations.org_id'), nullable=False)
    payload = Column(JSONB, nullable=False)
    metadata = Column(JSONB)
    risk_score = Column(Float)
    policy_results = Column(JSONB)
    trace_id = Column(String(255))
    parent_event_id = Column(UUID(as_uuid=True), ForeignKey('events.event_id'))
    estimated_cost = Column(Float)
    token_count = Column(Integer)
    created_at = Column(TIMESTAMP, default=datetime.utcnow)
    
    # Relationships
    agent = relationship("Agent", back_populates="events")
    risk_score_detail = relationship("RiskScore", back_populates="event", uselist=False)
```

### Alembic Migrations

#### Initial Migration

```python
"""Initial schema

Revision ID: 001
Create Date: 2024-01-01
"""

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

def upgrade():
    # Organizations
    op.create_table(
        'organizations',
        sa.Column('org_id', sa.String(255), primary_key=True),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('plan', sa.String(50), nullable=False),
        sa.Column('settings', postgresql.JSONB),
        sa.Column('created_at', sa.TIMESTAMP, server_default=sa.func.now()),
        sa.Column('updated_at', sa.TIMESTAMP, server_default=sa.func.now())
    )
    
    # Agents
    op.create_table(
        'agents',
        sa.Column('agent_id', sa.String(255), primary_key=True),
        sa.Column('org_id', sa.String(255), sa.ForeignKey('organizations.org_id')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('description', sa.Text),
        sa.Column('framework', sa.String(50)),
        sa.Column('config', postgresql.JSONB),
        sa.Column('status', sa.String(50), server_default='active'),
        sa.Column('created_at', sa.TIMESTAMP, server_default=sa.func.now()),
        sa.Column('updated_at', sa.TIMESTAMP, server_default=sa.func.now())
    )
    
    op.create_index('idx_agents_org_id', 'agents', ['org_id'])
    
    # Events
    op.create_table(
        'events',
        sa.Column('event_id', postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column('event_type', sa.String(50), nullable=False),
        sa.Column('timestamp', sa.TIMESTAMP, nullable=False),
        sa.Column('agent_id', sa.String(255), sa.ForeignKey('agents.agent_id')),
        sa.Column('org_id', sa.String(255), sa.ForeignKey('organizations.org_id')),
        sa.Column('payload', postgresql.JSONB, nullable=False),
        sa.Column('metadata', postgresql.JSONB),
        sa.Column('risk_score', sa.Float),
        sa.Column('policy_results', postgresql.JSONB),
        sa.Column('trace_id', sa.String(255)),
        sa.Column('parent_event_id', postgresql.UUID(as_uuid=True)),
        sa.Column('estimated_cost', sa.Float),
        sa.Column('token_count', sa.Integer),
        sa.Column('created_at', sa.TIMESTAMP, server_default=sa.func.now())
    )
    
    op.create_index('idx_events_agent_id', 'events', ['agent_id'])
    op.create_index('idx_events_org_id', 'events', ['org_id'])
    op.create_index('idx_events_timestamp', 'events', ['timestamp'])
    op.create_index('idx_events_type', 'events', ['event_type'])

def downgrade():
    op.drop_table('events')
    op.drop_table('agents')
    op.drop_table('organizations')
```

---

## Storage Strategy

### Hot vs Cold Storage

#### Hot Storage (PostgreSQL)
- Recent events (last 90 days)
- Active agents and policies
- Real-time queries
- Dashboard data

#### Cold Storage (S3/Archive)
- Historical events (>90 days)
- Compliance archives
- Batch analytics
- Cost-effective long-term retention

### Query Optimization

#### Materialized Views

```sql
-- Agent activity summary
CREATE MATERIALIZED VIEW agent_activity_summary AS
SELECT 
    agent_id,
    org_id,
    DATE(timestamp) as date,
    COUNT(*) as event_count,
    SUM(estimated_cost) as total_cost,
    SUM(token_count) as total_tokens,
    AVG(risk_score) as avg_risk_score,
    COUNT(*) FILTER (WHERE risk_score > 0.7) as high_risk_count
FROM events
WHERE timestamp > NOW() - INTERVAL '30 days'
GROUP BY agent_id, org_id, DATE(timestamp);

CREATE INDEX idx_agent_summary_agent ON agent_activity_summary(agent_id);
CREATE INDEX idx_agent_summary_date ON agent_activity_summary(date DESC);

-- Refresh periodically
REFRESH MATERIALIZED VIEW CONCURRENTLY agent_activity_summary;
```

#### Aggregation Tables

```sql
-- Pre-aggregate daily metrics
CREATE TABLE daily_metrics (
    metric_id UUID PRIMARY KEY,
    org_id VARCHAR(255) NOT NULL,
    agent_id VARCHAR(255),
    date DATE NOT NULL,
    event_count INTEGER,
    total_cost FLOAT,
    total_tokens INTEGER,
    avg_risk_score FLOAT,
    high_risk_count INTEGER,
    policy_violation_count INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_daily_metrics_unique 
ON daily_metrics(org_id, agent_id, date);
```

---

## Scalability and Performance

### V1 Performance Targets

- **Event Ingestion**: 100 events/minute
- **API Response Time**: <200ms (p95)
- **Dashboard Load Time**: <2 seconds
- **Database Queries**: <100ms (p95)
- **SDK Overhead**: <50ms per wrapped call

### Scaling Strategy

#### Vertical Scaling (V1)
- Single PostgreSQL instance
- Increase CPU/RAM as needed
- Connection pooling
- Query optimization

#### Horizontal Scaling (V2+)

```
┌─────────────────────────────────────────────────────────┐
│                    Load Balancer                         │
└────────────┬────────────────────────────────────────────┘
             │
             ├──────────┬──────────┬──────────┐
             ▼          ▼          ▼          ▼
        ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
        │ API    │ │ API    │ │ API    │ │ API    │
        │ Server │ │ Server │ │ Server │ │ Server │
        └────┬───┘ └────┬───┘ └────┬───┘ └────┬───┘
             │          │          │          │
             └──────────┴──────────┴──────────┘
                        │
                        ▼
             ┌─────────────────────┐
             │   Message Queue     │
             │   (Redis/Kafka)     │
             └──────────┬──────────┘
                        │
             ┌──────────┴──────────┐
             ▼                     ▼
        ┌─────────┐          ┌─────────┐
        │ Worker  │          │ Worker  │
        │ Pool    │          │ Pool    │
        └────┬────┘          └────┬────┘
             │                    │
             └────────┬───────────┘
                      ▼
             ┌─────────────────┐
             │   PostgreSQL    │
             │   (Primary)     │
             └────────┬────────┘
                      │
             ┌────────┴────────┐
             ▼                 ▼
        ┌─────────┐       ┌─────────┐
        │ Replica │       │ Replica │
        └─────────┘       └─────────┘
```

### Caching Strategy

#### Redis Cache Layers

```python
class CacheManager:
    """Multi-layer caching strategy"""
    
    def __init__(self, redis_client):
        self.redis = redis_client
        self.ttl = {
            'agent_config': 3600,      # 1 hour
            'policies': 1800,          # 30 minutes
            'risk_scores': 300,        # 5 minutes
            'dashboard_data': 60       # 1 minute
        }
    
    async def get_agent_config(self, agent_id: str):
        """Get agent config with caching"""
        cache_key = f"agent:config:{agent_id}"
        
        # Try cache first
        cached = await self.redis.get(cache_key)
        if cached:
            return json.loads(cached)
        
        # Fetch from database
        config = await db.get_agent_config(agent_id)
        
        # Cache for future requests
        await self.redis.setex(
            cache_key,
            self.ttl['agent_config'],
            json.dumps(config)
        )
        
        return config
```

### Database Connection Pooling

```python
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

# Connection pool configuration
engine = create_async_engine(
    DATABASE_URL,
    pool_size=20,           # Number of connections to maintain
    max_overflow=10,        # Additional connections when pool is full
    pool_timeout=30,        # Timeout waiting for connection
    pool_recycle=3600,      # Recycle connections after 1 hour
    pool_pre_ping=True      # Verify connections before use
)

AsyncSessionLocal = sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False
)
```

### Async Processing

```python
import asyncio
from typing import List

class AsyncEventProcessor:
    """Process events asynchronously"""
    
    async def process_batch(self, events: List[Event]):
        """Process multiple events concurrently"""
        tasks = [self.process_event(event) for event in events]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        return results
    
    async def process_event(self, event: Event):
        """Process single event"""
        # Calculate risk score
        risk_task = asyncio.create_task(self.calculate_risk(event))
        
        # Evaluate policies
        policy_task = asyncio.create_task(self.evaluate_policies(event))
        
        # Wait for both
        risk_score, policy_results = await asyncio.gather(
            risk_task,
            policy_task
        )
        
        # Persist to database
        await self.save_event(event, risk_score, policy_results)
```

---

## Multi-Tenancy Strategy

### V1: Single Tenant

Each customer gets their own deployment:
- Separate database instance
- Isolated infrastructure
- Full data isolation
- Simple security model

### V2+: Multi-Tenant Architecture

#### Row-Level Multi-Tenancy

All customers share infrastructure but data is isolated by org_id:

```sql
-- Enable row-level security
ALTER TABLE events ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see their org's data
CREATE POLICY org_isolation ON events
    USING (org_id = current_setting('app.current_org_id')::text);

-- Set org_id for session
SET app.current_org_id = 'org_123';
```

#### Schema-Based Multi-Tenancy

Each customer gets their own schema:

```sql
-- Create schema per organization
CREATE SCHEMA org_123;
CREATE SCHEMA org_456;

-- Create tables in each schema
CREATE TABLE org_123.events (...);
CREATE TABLE org_456.events (...);

-- Route queries to correct schema
SET search_path TO org_123;
```

#### Database-Per-Tenant

Largest customers get dedicated databases:

```python
class TenantRouter:
    """Route queries to correct database"""
    
    def __init__(self):
        self.connections = {
            'org_enterprise_1': 'postgresql://db1.example.com',
            'org_enterprise_2': 'postgresql://db2.example.com',
            'shared': 'postgresql://shared.example.com'
        }
    
    def get_connection(self, org_id: str):
        """Get database connection for org"""
        if org_id in self.connections:
            return self.connections[org_id]
        return self.connections['shared']
```

### Tenant Isolation

```python
class TenantContext:
    """Manage tenant context throughout request"""
    
    def __init__(self):
        self._context = contextvars.ContextVar('tenant_context')
    
    def set_org_id(self, org_id: str):
        """Set current organization"""
        self._context.set({'org_id': org_id})
    
    def get_org_id(self) -> str:
        """Get current organization"""
        context = self._context.get()
        return context['org_id']

# Middleware to set tenant context
@app.middleware("http")
async def tenant_middleware(request: Request, call_next):
    # Extract org_id from API key or JWT
    org_id = extract_org_id(request)
    
    # Set tenant context
    tenant_context.set_org_id(org_id)
    
    # Process request
    response = await call_next(request)
    
    return response
```

---

## Security Model

### Authentication

#### API Key Authentication

```python
from fastapi import Security, HTTPException
from fastapi.security import APIKeyHeader

api_key_header = APIKeyHeader(name="X-API-Key")

async def verify_api_key(api_key: str = Security(api_key_header)):
    """Verify API key and return org context"""
    # Hash the provided key
    key_hash = hash_api_key(api_key)
    
    # Look up in database
    key_record = await db.get_api_key(key_hash)
    
    if not key_record or key_record.revoked_at:
        raise HTTPException(status_code=401, detail="Invalid API key")
    
    if key_record.expires_at and key_record.expires_at < datetime.utcnow():
        raise HTTPException(status_code=401, detail="API key expired")
    
    # Update last used timestamp
    await db.update_api_key_last_used(key_record.key_id)
    
    return {
        'org_id': key_record.org_id,
        'permissions': key_record.permissions
    }
```

#### JWT Authentication (Dashboard)

```python
from jose import JWTError, jwt
from datetime import datetime, timedelta

SECRET_KEY = os.getenv("JWT_SECRET_KEY")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 60

def create_access_token(data: dict):
    """Create JWT access token"""
    to_encode = data.copy()
    expire = datetime.utcnow() + timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

async def verify_token(token: str):
    """Verify JWT token"""
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        user_id: str = payload.get("sub")
        org_id: str = payload.get("org_id")
        if user_id is None or org_id is None:
            raise HTTPException(status_code=401, detail="Invalid token")
        return {"user_id": user_id, "org_id": org_id}
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")
```

### Authorization

#### Role-Based Access Control (RBAC)

```python
from enum import Enum

class Role(str, Enum):
    ADMIN = "admin"
    DEVELOPER = "developer"
    VIEWER = "viewer"

class Permission(str, Enum):
    READ_EVENTS = "read:events"
    WRITE_EVENTS = "write:events"
    READ_POLICIES = "read:policies"
    WRITE_POLICIES = "write:policies"
    MANAGE_AGENTS = "manage:agents"
    MANAGE_USERS = "manage:users"

ROLE_PERMISSIONS = {
    Role.ADMIN: [
        Permission.READ_EVENTS,
        Permission.WRITE_EVENTS,
        Permission.READ_POLICIES,
        Permission.WRITE_POLICIES,
        Permission.MANAGE_AGENTS,
        Permission.MANAGE_USERS
    ],
    Role.DEVELOPER: [
        Permission.READ_EVENTS,
        Permission.WRITE_EVENTS,
        Permission.READ_POLICIES,
        Permission.MANAGE_AGENTS
    ],
    Role.VIEWER: [
        Permission.READ_EVENTS,
        Permission.READ_POLICIES
    ]
}

def require_permission(permission: Permission):
    """Decorator to require specific permission"""
    def decorator(func):
        async def wrapper(*args, current_user: dict, **kwargs):
            user_role = current_user.get('role')
            user_permissions = ROLE_PERMISSIONS.get(user_role, [])
            
            if permission not in user_permissions:
                raise HTTPException(
                    status_code=403,
                    detail=f"Permission denied: {permission}"
                )
            
            return await func(*args, current_user=current_user, **kwargs)
        return wrapper
    return decorator
```

### Data Encryption

#### At Rest

```python
from cryptography.fernet import Fernet

class EncryptionService:
    """Encrypt sensitive data at rest"""
    
    def __init__(self, key: bytes):
        self.cipher = Fernet(key)
    
    def encrypt(self, data: str) -> str:
        """Encrypt string data"""
        encrypted = self.cipher.encrypt(data.encode())
        return encrypted.decode()
    
    def decrypt(self, encrypted_data: str) -> str:
        """Decrypt string data"""
        decrypted = self.cipher.decrypt(encrypted_data.encode())
        return decrypted.decode()

# Encrypt sensitive fields before storage
class SensitiveEvent(Event):
    def save(self):
        # Encrypt PII in payload
        if 'email' in self.payload:
            self.payload['email'] = encryption_service.encrypt(
                self.payload['email']
            )
        super().save()
```

#### In Transit

```python
# Force HTTPS in production
from fastapi.middleware.httpsredirect import HTTPSRedirectMiddleware

if os.getenv("ENVIRONMENT") == "production":
    app.add_middleware(HTTPSRedirectMiddleware)

# TLS configuration for database connections
DATABASE_URL = "postgresql://user:pass@host:5432/db?sslmode=require"
```

### Audit Logging

```python
class AuditLogger:
    """Log security-relevant events"""
    
    async def log_access(
        self,
        user_id: str,
        action: str,
        resource: str,
        result: str
    ):
        """Log access attempt"""
        await db.insert_audit_log({
            'user_id': user_id,
            'action': action,
            'resource': resource,
            'result': result,  # success, denied, error
            'timestamp': datetime.utcnow(),
            'ip_address': get_client_ip(),
            'user_agent': get_user_agent()
        })

# Usage
@app.post("/api/v1/policies")
async def create_policy(
    policy: Policy,
    current_user: dict = Depends(verify_token)
):
    try:
        result = await policy_service.create(policy)
        await audit_logger.log_access(
            user_id=current_user['user_id'],
            action='CREATE_POLICY',
            resource=f"policy:{policy.policy_id}",
            result='success'
        )
        return result
    except Exception as e:
        await audit_logger.log_access(
            user_id=current_user['user_id'],
            action='CREATE_POLICY',
            resource=f"policy:{policy.policy_id}",
            result='error'
        )
        raise
```

---

## Observability and Monitoring

### Structured Logging

```python
from loguru import logger
import sys

# Configure Loguru
logger.remove()  # Remove default handler
logger.add(
    sys.stdout,
    format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {name}:{function}:{line} | {message}",
    level="INFO",
    serialize=True  # JSON output
)

# Add file handler
logger.add(
    "logs/app.log",
    rotation="500 MB",
    retention="10 days",
    compression="zip"
)

# Usage with context
logger.info(
    "Event processed",
    event_id=event.event_id,
    agent_id=event.agent_id,
    risk_score=event.risk_score
)
```

### Metrics Collection

```python
from prometheus_client import Counter, Histogram, Gauge, generate_latest

# Define metrics
events_processed = Counter(
    'events_processed_total',
    'Total number of events processed',
    ['event_type', 'agent_id']
)

event_processing_time = Histogram(
    'event_processing_seconds',
    'Time spent processing events',
    ['event_type']
)

active_agents = Gauge(
    'active_agents',
    'Number of active agents',
    ['org_id']
)

risk_score_distribution = Histogram(
    'risk_score_distribution',
    'Distribution of risk scores',
    buckets=[0.1, 0.3, 0.5, 0.7, 0.9, 1.0]
)

# Usage
@app.post("/api/v1/events")
async def create_event(event: Event):
    with event_processing_time.labels(event.event_type).time():
        result = await process_event(event)
        events_processed.labels(
            event_type=event.event_type,
            agent_id=event.agent_id
        ).inc()
        risk_score_distribution.observe(result.risk_score)
        return result

# Expose metrics endpoint
@app.get("/metrics")
async def metrics():
    return Response(
        content=generate_latest(),
        media_type="text/plain"
    )
```

### Distributed Tracing

```python
import uuid
from contextvars import ContextVar

trace_id_var = ContextVar('trace_id', default=None)

@app.middleware("http")
async def tracing_middleware(request: Request, call_next):
    # Generate or extract trace ID
    trace_id = request.headers.get('X-Trace-ID') or str(uuid.uuid4())
    trace_id_var.set(trace_id)
    
    # Add to response headers
    response = await call_next(request)
    response.headers['X-Trace-ID'] = trace_id
    
    return response

# Use in logging
def log_with_trace(message: str, **kwargs):
    trace_id = trace_id_var.get()
    logger.info(message, trace_id=trace_id, **kwargs)
```

### Health Checks

```python
@app.get("/health")
async def health_check():
    """Basic health check"""
    return {"status": "healthy"}

@app.get("/health/detailed")
async def detailed_health_check():
    """Detailed health check with dependencies"""
    health = {
        "status": "healthy",
        "checks": {}
    }
    
    # Check database
    try:
        await db.execute("SELECT 1")
        health["checks"]["database"] = "healthy"
    except Exception as e:
        health["checks"]["database"] = f"unhealthy: {str(e)}"
        health["status"] = "unhealthy"
    
    # Check Redis
    try:
        await redis.ping()
        health["checks"]["redis"] = "healthy"
    except Exception as e:
        health["checks"]["redis"] = f"unhealthy: {str(e)}"
        health["status"] = "degraded"
    
    return health
```

---


## API Design

### REST API Endpoints

#### Events API

```python
# Create event
POST /api/v1/events
Content-Type: application/json
X-API-Key: your_api_key

{
  "event_type": "PROMPT_SENT",
  "agent_id": "agent_123",
  "payload": {...},
  "metadata": {...}
}

Response: 201 Created
{
  "event_id": "uuid",
  "risk_score": 0.3,
  "policy_results": [...],
  "action": "ALLOW"
}

# Get events
GET /api/v1/events?agent_id=agent_123&limit=100&offset=0

Response: 200 OK
{
  "events": [...],
  "total": 1000,
  "limit": 100,
  "offset": 0
}

# Get single event
GET /api/v1/events/{event_id}

Response: 200 OK
{
  "event_id": "uuid",
  "event_type": "PROMPT_SENT",
  ...
}
```

#### Agents API

```python
# List agents
GET /api/v1/agents

Response: 200 OK
{
  "agents": [
    {
      "agent_id": "agent_123",
      "name": "Customer Support Agent",
      "status": "active",
      "framework": "langchain"
    }
  ]
}

# Create agent
POST /api/v1/agents
{
  "agent_id": "agent_123",
  "name": "Customer Support Agent",
  "framework": "langchain",
  "config": {...}
}

# Update agent
PUT /api/v1/agents/{agent_id}
{
  "name": "Updated Name",
  "status": "paused"
}

# Delete agent
DELETE /api/v1/agents/{agent_id}
```

#### Policies API

```python
# List policies
GET /api/v1/policies

# Create policy
POST /api/v1/policies
{
  "policy_id": "refund-limit",
  "name": "Refund Limit Policy",
  "policy_type": "approval",
  "config": {...}
}

# Evaluate policy (test)
POST /api/v1/policies/evaluate
{
  "event": {...},
  "policy_ids": ["policy_1", "policy_2"]
}

Response: 200 OK
{
  "action": "BLOCK",
  "policy_results": [...]
}
```

#### Analytics API

```python
# Get agent metrics
GET /api/v1/analytics/agents/{agent_id}/metrics?start_date=2024-01-01&end_date=2024-01-31

Response: 200 OK
{
  "event_count": 10000,
  "total_cost": 150.50,
  "total_tokens": 500000,
  "avg_risk_score": 0.25,
  "high_risk_count": 50
}

# Get risk distribution
GET /api/v1/analytics/risk-distribution?org_id=org_123

Response: 200 OK
{
  "low": 8500,
  "medium": 1200,
  "high": 300
}
```

### FastAPI Implementation

```python
from fastapi import FastAPI, Depends, HTTPException, Query
from typing import List, Optional

app = FastAPI(
    title="AI Control Layer API",
    version="1.0.0",
    description="API for AI agent governance and control"
)

# Events endpoints
@app.post("/api/v1/events", status_code=201)
async def create_event(
    event: EventCreate,
    auth: dict = Depends(verify_api_key)
):
    """Create and process a new event"""
    # Set org_id from auth
    event.org_id = auth['org_id']
    
    # Calculate risk score
    risk_score = await risk_engine.calculate_risk(event)
    
    # Evaluate policies
    policy_result = await policy_engine.evaluate(event)
    
    # Save event
    saved_event = await event_service.create(
        event,
        risk_score=risk_score,
        policy_results=policy_result.policy_results
    )
    
    return {
        "event_id": saved_event.event_id,
        "risk_score": risk_score.composite,
        "policy_results": policy_result.policy_results,
        "action": policy_result.action
    }

@app.get("/api/v1/events")
async def list_events(
    agent_id: Optional[str] = None,
    event_type: Optional[str] = None,
    start_date: Optional[datetime] = None,
    end_date: Optional[datetime] = None,
    limit: int = Query(100, le=1000),
    offset: int = Query(0, ge=0),
    auth: dict = Depends(verify_api_key)
):
    """List events with filtering"""
    events = await event_service.list(
        org_id=auth['org_id'],
        agent_id=agent_id,
        event_type=event_type,
        start_date=start_date,
        end_date=end_date,
        limit=limit,
        offset=offset
    )
    
    total = await event_service.count(
        org_id=auth['org_id'],
        agent_id=agent_id,
        event_type=event_type,
        start_date=start_date,
        end_date=end_date
    )
    
    return {
        "events": events,
        "total": total,
        "limit": limit,
        "offset": offset
    }
```

### Error Handling

```python
from fastapi import Request
from fastapi.responses import JSONResponse

class APIError(Exception):
    def __init__(self, status_code: int, message: str, details: dict = None):
        self.status_code = status_code
        self.message = message
        self.details = details or {}

@app.exception_handler(APIError)
async def api_error_handler(request: Request, exc: APIError):
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "error": {
                "message": exc.message,
                "details": exc.details,
                "trace_id": trace_id_var.get()
            }
        }
    )

@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={
            "error": {
                "message": "Internal server error",
                "trace_id": trace_id_var.get()
            }
        }
    )
```

### Rate Limiting

```python
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address
from slowapi.errors import RateLimitExceeded

limiter = Limiter(key_func=get_remote_address)
app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)

@app.post("/api/v1/events")
@limiter.limit("100/minute")
async def create_event(request: Request, event: EventCreate):
    # ... implementation
    pass
```

---

## Frontend Dashboard

### Technology Stack

- **React 18**: Component-based UI
- **Vite**: Fast build tool and dev server
- **TypeScript**: Type safety
- **TanStack Query**: Data fetching and caching
- **Tailwind CSS**: Utility-first styling
- **Recharts**: Data visualization
- **React Router**: Client-side routing

### Project Structure

```
dashboard/
├── src/
│   ├── components/
│   │   ├── EventList.tsx
│   │   ├── AgentCard.tsx
│   │   ├── RiskChart.tsx
│   │   └── PolicyManager.tsx
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── Events.tsx
│   │   ├── Agents.tsx
│   │   └── Policies.tsx
│   ├── services/
│   │   ├── api.ts
│   │   └── auth.ts
│   ├── hooks/
│   │   ├── useEvents.ts
│   │   └── useAgents.ts
│   ├── types/
│   │   └── index.ts
│   ├── App.tsx
│   └── main.tsx
├── package.json
└── vite.config.ts
```

### Key Components

#### Event List Component

```typescript
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getEvents } from '../services/api';

interface Event {
  event_id: string;
  event_type: string;
  timestamp: string;
  agent_id: string;
  risk_score: number;
}

export const EventList: React.FC = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['events'],
    queryFn: () => getEvents({ limit: 100 }),
    refetchInterval: 5000 // Refresh every 5 seconds
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading events</div>;

  return (
    <div className="space-y-4">
      {data?.events.map((event: Event) => (
        <div key={event.event_id} className="border rounded p-4">
          <div className="flex justify-between">
            <span className="font-semibold">{event.event_type}</span>
            <span className={`px-2 py-1 rounded ${
              event.risk_score > 0.7 ? 'bg-red-100 text-red-800' :
              event.risk_score > 0.3 ? 'bg-yellow-100 text-yellow-800' :
              'bg-green-100 text-green-800'
            }`}>
              Risk: {(event.risk_score * 100).toFixed(0)}%
            </span>
          </div>
          <div className="text-sm text-gray-600 mt-2">
            Agent: {event.agent_id} | {new Date(event.timestamp).toLocaleString()}
          </div>
        </div>
      ))}
    </div>
  );
};
```

#### Risk Chart Component

```typescript
import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import { useQuery } from '@tanstack/react-query';
import { getAgentMetrics } from '../services/api';

export const RiskChart: React.FC<{ agentId: string }> = ({ agentId }) => {
  const { data } = useQuery({
    queryKey: ['metrics', agentId],
    queryFn: () => getAgentMetrics(agentId, { days: 7 })
  });

  return (
    <LineChart width={600} height={300} data={data?.daily_metrics}>
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="date" />
      <YAxis />
      <Tooltip />
      <Legend />
      <Line type="monotone" dataKey="avg_risk_score" stroke="#8884d8" name="Avg Risk Score" />
      <Line type="monotone" dataKey="high_risk_count" stroke="#ff0000" name="High Risk Events" />
    </LineChart>
  );
};
```

#### API Service

```typescript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';
const API_KEY = localStorage.getItem('api_key');

async function fetchAPI(endpoint: string, options: RequestInit = {}) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': API_KEY || '',
      ...options.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
}

export const getEvents = (params: {
  agent_id?: string;
  limit?: number;
  offset?: number;
}) => {
  const queryString = new URLSearchParams(params as any).toString();
  return fetchAPI(`/api/v1/events?${queryString}`);
};

export const getAgents = () => {
  return fetchAPI('/api/v1/agents');
};

export const getPolicies = () => {
  return fetchAPI('/api/v1/policies');
};

export const getAgentMetrics = (agentId: string, params: { days: number }) => {
  return fetchAPI(`/api/v1/analytics/agents/${agentId}/metrics?days=${params.days}`);
};
```

### Dashboard Views

#### Main Dashboard

```typescript
import React from 'react';
import { EventList } from '../components/EventList';
import { RiskChart } from '../components/RiskChart';
import { AgentCard } from '../components/AgentCard';

export const Dashboard: React.FC = () => {
  return (
    <div className="container mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">AI Control Layer Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-2">Total Events</h3>
          <p className="text-3xl font-bold">10,234</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-2">High Risk Events</h3>
          <p className="text-3xl font-bold text-red-600">23</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-2">Active Agents</h3>
          <p className="text-3xl font-bold">5</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Recent Events</h2>
          <EventList />
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Risk Trends</h2>
          <RiskChart agentId="all" />
        </div>
      </div>
    </div>
  );
};
```

---

## Testing Strategy

### Unit Tests

```python
import pytest
from ai_control.risk_engine import RiskEngine, detect_hallucination

def test_hallucination_detection():
    """Test hallucination detection heuristics"""
    # High confidence, no hedging
    score = detect_hallucination("This is a definitive answer.", confidence=0.95)
    assert score < 0.3
    
    # Low confidence, hedging language
    score = detect_hallucination("I think maybe this could be the answer.", confidence=0.5)
    assert score > 0.5
    
    # Very short response
    score = detect_hallucination("Yes.", confidence=0.8)
    assert score > 0.2

def test_risk_engine_composite_score():
    """Test composite risk scoring"""
    engine = RiskEngine()
    event = Event(
        event_type="RESPONSE_RECEIVED",
        agent_id="test_agent",
        org_id="test_org",
        payload={"response": "Test response", "confidence": 0.9}
    )
    
    risk_score = engine.calculate_risk(event, context={})
    assert 0.0 <= risk_score.composite <= 1.0
    assert risk_score.severity in ["LOW", "MEDIUM", "HIGH"]
```

### Integration Tests

```python
import pytest
from httpx import AsyncClient
from main import app

@pytest.mark.asyncio
async def test_create_event():
    """Test event creation endpoint"""
    async with AsyncClient(app=app, base_url="http://test") as client:
        response = await client.post(
            "/api/v1/events",
            json={
                "event_type": "PROMPT_SENT",
                "agent_id": "test_agent",
                "payload": {"prompt": "Test prompt"}
            },
            headers={"X-API-Key": "test_key"}
        )
        
        assert response.status_code == 201
        data = response.json()
        assert "event_id" in data
        assert "risk_score" in data

@pytest.mark.asyncio
async def test_policy_enforcement():
    """Test policy enforcement"""
    # Create a blocking policy
    policy = Policy(
        policy_id="test_policy",
        policy_type="access_control",
        config={"denied_actions": ["dangerous_action"]}
    )
    await policy_service.create(policy)
    
    # Try to perform denied action
    event = Event(
        event_type="TOOL_CALLED",
        agent_id="test_agent",
        payload={"action": "dangerous_action"}
    )
    
    result = await policy_engine.evaluate(event)
    assert result.action == "BLOCK"
```

### Load Tests

```python
import asyncio
from locust import HttpUser, task, between

class AIControlUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        """Set up API key"""
        self.headers = {"X-API-Key": "test_key"}
    
    @task(3)
    def create_event(self):
        """Simulate event creation"""
        self.client.post(
            "/api/v1/events",
            json={
                "event_type": "PROMPT_SENT",
                "agent_id": "load_test_agent",
                "payload": {"prompt": "Test prompt"}
            },
            headers=self.headers
        )
    
    @task(1)
    def list_events(self):
        """Simulate event listing"""
        self.client.get(
            "/api/v1/events?limit=100",
            headers=self.headers
        )

# Run: locust -f load_test.py --host=http://localhost:8000
```

### End-to-End Tests

```python
import pytest
from ai_control import ControlLayer

@pytest.mark.e2e
async def test_full_agent_workflow():
    """Test complete agent workflow with SDK"""
    # Initialize SDK
    control = ControlLayer(
        api_key="test_key",
        org_id="test_org"
    )
    
    # Wrap agent function
    @control.monitor(agent_id="e2e_test_agent")
    def process_request(request):
        return f"Processed: {request}"
    
    # Execute function
    result = process_request("test request")
    assert result == "Processed: test request"
    
    # Verify event was logged
    events = await control.get_events(agent_id="e2e_test_agent")
    assert len(events) > 0
    assert events[0].event_type == "FUNCTION_START"
```

---

## Data Science Integration

### Feature Engineering

```python
class FeatureExtractor:
    """Extract features from events for ML models"""
    
    def extract(self, event: Event) -> dict:
        """Extract features from event"""
        response = event.payload.get('response', '')
        
        features = {
            # Text features
            'response_length': len(response),
            'word_count': len(response.split()),
            'avg_word_length': np.mean([len(w) for w in response.split()]),
            'sentence_count': response.count('.') + response.count('!') + response.count('?'),
            
            # Linguistic features
            'hedge_word_count': sum(1 for word in ['maybe', 'perhaps', 'possibly'] if word in response.lower()),
            'caps_ratio': sum(1 for c in response if c.isupper()) / max(len(response), 1),
            'punctuation_ratio': sum(1 for c in response if c in '!?.,;:') / max(len(response), 1),
            
            # Metadata features
            'confidence': event.payload.get('confidence', 1.0),
            'token_count': event.token_count or 0,
            'estimated_cost': event.estimated_cost or 0.0,
            
            # Temporal features
            'hour_of_day': event.timestamp.hour,
            'day_of_week': event.timestamp.weekday(),
            
            # Agent features
            'agent_id_hash': hash(event.agent_id) % 1000
        }
        
        return features
```

### Model Training Pipeline

```python
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
import joblib

class RiskModelTrainer:
    """Train ML models for risk prediction"""
    
    def __init__(self):
        self.feature_extractor = FeatureExtractor()
        self.model = RandomForestClassifier(n_estimators=100)
    
    async def collect_training_data(self, days: int = 30):
        """Collect labeled training data"""
        # Get events with human-reviewed risk scores
        events = await db.get_events_with_labels(
            start_date=datetime.utcnow() - timedelta(days=days)
        )
        
        X = []
        y = []
        
        for event in events:
            features = self.feature_extractor.extract(event)
            X.append(list(features.values()))
            y.append(event.human_risk_score)  # Human-labeled score
        
        return np.array(X), np.array(y)
    
    def train(self, X, y):
        """Train the model"""
        X_train, X_test, y_train, y_test = train_test_split(
            X, y, test_size=0.2, random_state=42
        )
        
        self.model.fit(X_train, y_train)
        
        # Evaluate
        train_score = self.model.score(X_train, y_train)
        test_score = self.model.score(X_test, y_test)
        
        print(f"Train score: {train_score:.3f}")
        print(f"Test score: {test_score:.3f}")
        
        return test_score
    
    def save(self, path: str):
        """Save trained model"""
        joblib.dump(self.model, path)
    
    def load(self, path: str):
        """Load trained model"""
        self.model = joblib.load(path)
```

### Offline Analysis

```python
import pandas as pd
import matplotlib.pyplot as plt

class AgentAnalyzer:
    """Analyze agent behavior patterns"""
    
    async def analyze_agent(self, agent_id: str, days: int = 30):
        """Generate analysis report for agent"""
        events = await db.get_events(
            agent_id=agent_id,
            start_date=datetime.utcnow() - timedelta(days=days)
        )
        
        df = pd.DataFrame([e.dict() for e in events])
        
        analysis = {
            'total_events': len(df),
            'avg_risk_score': df['risk_score'].mean(),
            'high_risk_events': len(df[df['risk_score'] > 0.7]),
            'total_cost': df['estimated_cost'].sum(),
            'total_tokens': df['token_count'].sum(),
            'event_types': df['event_type'].value_counts().to_dict(),
            'risk_trend': df.groupby(df['timestamp'].dt.date)['risk_score'].mean().to_dict()
        }
        
        return analysis
    
    def plot_risk_distribution(self, events: List[Event]):
        """Plot risk score distribution"""
        risk_scores = [e.risk_score for e in events if e.risk_score]
        
        plt.figure(figsize=(10, 6))
        plt.hist(risk_scores, bins=20, edgecolor='black')
        plt.xlabel('Risk Score')
        plt.ylabel('Frequency')
        plt.title('Risk Score Distribution')
        plt.savefig('risk_distribution.png')
```

---


## Failure Scenarios and Resilience

### Failure Modes

#### 1. SDK Failure

**Scenario**: SDK crashes or encounters error

**Impact**: Agent operations disrupted

**Mitigation**:
```python
class ResilientControlLayer(ControlLayer):
    """SDK with failure resilience"""
    
    def __init__(self, *args, fail_open: bool = True, **kwargs):
        super().__init__(*args, **kwargs)
        self.fail_open = fail_open
        self.circuit_breaker = CircuitBreaker()
    
    def monitor(self, agent_id: str):
        def decorator(func):
            @wraps(func)
            def wrapper(*args, **kwargs):
                try:
                    # Try to log event
                    self.circuit_breaker.call(
                        self._log_event,
                        event_type="FUNCTION_START",
                        agent_id=agent_id
                    )
                except Exception as e:
                    logger.error(f"SDK error: {e}")
                    if not self.fail_open:
                        raise
                
                # Always execute agent function
                return func(*args, **kwargs)
            
            return wrapper
        return decorator
```

#### 2. Database Unavailable

**Scenario**: PostgreSQL connection lost

**Impact**: Cannot persist events

**Mitigation**:
```python
class EventBuffer:
    """Buffer events when database is unavailable"""
    
    def __init__(self, max_size: int = 10000):
        self.buffer = []
        self.max_size = max_size
    
    def add(self, event: Event):
        """Add event to buffer"""
        if len(self.buffer) < self.max_size:
            self.buffer.append(event)
        else:
            logger.warning("Event buffer full, dropping event")
    
    async def flush(self):
        """Flush buffered events to database"""
        if not self.buffer:
            return
        
        try:
            await db.bulk_insert_events(self.buffer)
            logger.info(f"Flushed {len(self.buffer)} buffered events")
            self.buffer.clear()
        except Exception as e:
            logger.error(f"Failed to flush buffer: {e}")

# Periodic flush
async def flush_buffer_periodically():
    while True:
        await asyncio.sleep(60)  # Every minute
        await event_buffer.flush()
```

#### 3. API Overload

**Scenario**: Too many requests overwhelm API

**Impact**: Slow response times, timeouts

**Mitigation**:
```python
from fastapi_limiter import FastAPILimiter
from fastapi_limiter.depends import RateLimiter

# Initialize rate limiter
await FastAPILimiter.init(redis)

# Apply rate limits
@app.post("/api/v1/events")
@limiter.limit("100/minute")
async def create_event(
    request: Request,
    event: EventCreate,
    rate_limit: RateLimiter = Depends(RateLimiter(times=100, seconds=60))
):
    # ... implementation
    pass

# Queue for async processing
from asyncio import Queue

event_queue = Queue(maxsize=10000)

@app.post("/api/v1/events")
async def create_event(event: EventCreate):
    """Accept event and queue for processing"""
    try:
        event_queue.put_nowait(event)
        return {"status": "queued", "event_id": event.event_id}
    except asyncio.QueueFull:
        raise HTTPException(status_code=503, detail="Service overloaded")

# Background worker
async def process_event_queue():
    while True:
        event = await event_queue.get()
        try:
            await process_event(event)
        except Exception as e:
            logger.error(f"Failed to process event: {e}")
        finally:
            event_queue.task_done()
```

#### 4. Agent Misbehavior

**Scenario**: Agent enters infinite loop or makes excessive calls

**Impact**: Cost overruns, resource exhaustion

**Mitigation**:
```python
class AgentThrottler:
    """Throttle agent actions to prevent abuse"""
    
    def __init__(self, redis_client):
        self.redis = redis_client
    
    async def check_rate_limit(self, agent_id: str, limit: int = 100, window: int = 60):
        """Check if agent exceeds rate limit"""
        key = f"rate_limit:{agent_id}"
        
        # Increment counter
        count = await self.redis.incr(key)
        
        # Set expiry on first request
        if count == 1:
            await self.redis.expire(key, window)
        
        if count > limit:
            raise RateLimitExceeded(
                f"Agent {agent_id} exceeded rate limit: {count}/{limit} in {window}s"
            )
        
        return count

# Use in SDK
@control.monitor(agent_id="agent_123")
async def agent_action():
    await throttler.check_rate_limit("agent_123", limit=100, window=60)
    # ... perform action
```

### Disaster Recovery

#### Backup Strategy

```python
# Daily database backups
import subprocess
from datetime import datetime

def backup_database():
    """Create database backup"""
    timestamp = datetime.utcnow().strftime("%Y%m%d_%H%M%S")
    backup_file = f"backup_{timestamp}.sql"
    
    subprocess.run([
        "pg_dump",
        "-h", DB_HOST,
        "-U", DB_USER,
        "-d", DB_NAME,
        "-f", backup_file
    ])
    
    # Upload to S3
    s3_client.upload_file(
        backup_file,
        "ai-control-backups",
        f"backups/{backup_file}"
    )
    
    logger.info(f"Database backup created: {backup_file}")

# Schedule daily backups
from apscheduler.schedulers.asyncio import AsyncIOScheduler

scheduler = AsyncIOScheduler()
scheduler.add_job(backup_database, 'cron', hour=2)  # 2 AM daily
scheduler.start()
```

#### Point-in-Time Recovery

```sql
-- Enable WAL archiving for PITR
ALTER SYSTEM SET wal_level = replica;
ALTER SYSTEM SET archive_mode = on;
ALTER SYSTEM SET archive_command = 'cp %p /archive/%f';

-- Restore to specific point in time
pg_restore --dbname=ai_control --time="2024-01-15 14:30:00" backup.sql
```

---

## Technical Debt Management

### Documentation Standards

#### Code Documentation

```python
def calculate_risk_score(event: Event, context: dict) -> RiskScore:
    """
    Calculate composite risk score for an agent event.
    
    Args:
        event: The event to evaluate
        context: Additional context including policies, budgets, etc.
    
    Returns:
        RiskScore object with composite score and dimension breakdown
    
    Raises:
        ValueError: If event is missing required fields
    
    Example:
        >>> event = Event(event_type="PROMPT_SENT", ...)
        >>> score = calculate_risk_score(event, context={})
        >>> print(score.composite)
        0.35
    
    Notes:
        - Uses weighted average of multiple risk dimensions
        - Weights can be customized via RiskEngine initialization
        - Future versions will use ML models instead of heuristics
    
    Technical Debt:
        - TODO: Replace heuristics with ML model (v2)
        - TODO: Add caching for repeated calculations
        - FIXME: Hallucination detection needs improvement
    """
    pass
```

#### Architecture Decision Records (ADR)

```markdown
# ADR-001: Use PostgreSQL for Event Storage

## Status
Accepted

## Context
Need to choose database for storing agent events. Requirements:
- ACID compliance for audit trails
- JSON support for flexible payloads
- Good performance for time-series queries
- Mature ecosystem

## Decision
Use PostgreSQL with JSONB columns

## Consequences
Positive:
- Strong consistency guarantees
- Flexible schema with JSONB
- Excellent tooling and community support

Negative:
- May need partitioning for very large datasets
- Not optimized for pure time-series workloads

## Alternatives Considered
- MongoDB: Less mature for ACID transactions
- TimescaleDB: Overkill for V1 scale
- DynamoDB: Vendor lock-in concerns
```

### Technical Debt Tracking

```python
# tech_debt.py
"""
Technical Debt Registry

This file tracks known technical debt items with priority and effort estimates.
"""

TECH_DEBT_ITEMS = [
    {
        "id": "TD-001",
        "title": "Replace heuristic risk scoring with ML models",
        "description": "Current risk engine uses simple heuristics. Need ML models for better accuracy.",
        "priority": "HIGH",
        "effort": "LARGE",
        "target_version": "v2.0",
        "created": "2024-01-01",
        "owner": "data-science-team"
    },
    {
        "id": "TD-002",
        "title": "Implement async event processing",
        "description": "Events currently processed synchronously. Need async queue for scale.",
        "priority": "MEDIUM",
        "effort": "MEDIUM",
        "target_version": "v1.5",
        "created": "2024-01-01",
        "owner": "backend-team"
    },
    {
        "id": "TD-003",
        "title": "Add comprehensive integration tests",
        "description": "Need more integration tests for SDK + API interactions.",
        "priority": "MEDIUM",
        "effort": "SMALL",
        "target_version": "v1.2",
        "created": "2024-01-01",
        "owner": "qa-team"
    }
]
```

### Refactoring Guidelines

1. **One Path Principle**: Build one implementation path, design for replacement
2. **Incremental Refactoring**: Small, testable changes
3. **Test Coverage**: Maintain >80% coverage during refactors
4. **Backward Compatibility**: Version APIs, deprecate gracefully
5. **Documentation**: Update docs with every refactor

---

## Engineering Principles

### 1. Build One Path, Design for Replacement

Don't over-engineer for hypothetical futures. Build the simplest thing that works, but design interfaces that allow replacement.

```python
# Good: Simple implementation, clean interface
class RiskEngine:
    def calculate_risk(self, event: Event) -> RiskScore:
        """Calculate risk score - implementation can be swapped"""
        return self._heuristic_scoring(event)  # V1: heuristics
        # return self._ml_scoring(event)       # V2: ML model

# Bad: Tightly coupled, hard to replace
def calculate_risk_with_heuristics_and_ml_and_rules(event):
    # Complex logic mixing multiple approaches
    pass
```

### 2. Versioned Interfaces

All APIs and schemas should be versioned from day one.

```python
# API versioning
@app.post("/api/v1/events")  # Version in URL
async def create_event_v1(event: EventCreateV1):
    pass

@app.post("/api/v2/events")  # New version, old still works
async def create_event_v2(event: EventCreateV2):
    pass

# Schema versioning
class Event(BaseModel):
    schema_version: str = "1.0"
    # ... fields
```

### 3. Incremental Iteration

Ship small, working increments. Don't wait for perfection.

**V1 Milestones**:
- Week 1: Basic SDK wrapper
- Week 2: Event logging to database
- Week 3: Simple risk scoring
- Week 4: Policy enforcement
- Week 5: Basic dashboard
- Week 6: Documentation and examples

### 4. Test First, Refactor Safely

Write tests before refactoring. Maintain test coverage.

```python
# 1. Write test for desired behavior
def test_risk_calculation():
    event = create_test_event()
    score = calculate_risk(event)
    assert 0.0 <= score <= 1.0

# 2. Refactor implementation
def calculate_risk(event):
    # New implementation
    pass

# 3. Verify tests still pass
pytest test_risk_engine.py
```

### 5. Fail Fast, Fail Loudly

Detect errors early and make them visible.

```python
# Good: Validate early
def create_event(event: Event):
    if not event.agent_id:
        raise ValueError("agent_id is required")
    if not event.org_id:
        raise ValueError("org_id is required")
    # ... continue

# Bad: Silent failures
def create_event(event: Event):
    agent_id = event.agent_id or "unknown"  # Masks the problem
    # ... continue
```

---

## Business Defensibility

### Competitive Moats

#### 1. Data Network Effects

As more agents are monitored, we collect more data:
- Better risk models trained on diverse agent behaviors
- Benchmark data for comparing agent performance
- Industry-specific risk patterns
- Proprietary dataset that competitors can't replicate

**Strategy**: Aggregate anonymized data across customers to improve risk models for everyone.

#### 2. Integration Depth

Deep integrations with agent frameworks create switching costs:
- Custom adapters for LangChain, AutoGen, etc.
- Framework-specific optimizations
- Extensive testing and validation
- Documentation and examples

**Strategy**: Become the default governance layer for popular frameworks.

#### 3. Compliance Positioning

Become the standard for AI governance in regulated industries:
- SOC 2 Type II certification
- HIPAA compliance
- GDPR compliance
- Industry-specific certifications (finance, healthcare)

**Strategy**: Partner with compliance auditors and regulators.

#### 4. Ecosystem Lock-In

Build an ecosystem around the platform:
- Third-party integrations (Slack, PagerDuty, etc.)
- Custom policy templates marketplace
- Community-contributed risk models
- Training and certification programs

**Strategy**: Create a platform that becomes more valuable with ecosystem growth.

### Pricing Strategy

#### Tier Structure

**Free Tier**:
- 1 agent
- 10,000 events/month
- 7-day data retention
- Community support

**Pro Tier** ($99/month):
- 10 agents
- 100,000 events/month
- 90-day data retention
- Email support
- Basic analytics

**Enterprise Tier** (Custom):
- Unlimited agents
- Unlimited events
- Custom data retention
- Dedicated support
- Advanced analytics
- Custom integrations
- SLA guarantees

#### Value Metrics

Price based on value delivered:
- Number of agents monitored
- Event volume
- Data retention period
- Support level
- Advanced features (ML models, custom policies)

---

## Execution Constraints

### Founder Time Constraints

**Available Time**: 1 hour per day

**Optimization Strategy**:
1. Focus on high-leverage activities
2. Automate repetitive tasks
3. Use AI coding assistants
4. Prioritize ruthlessly
5. Build in public for accountability

### Daily Workflow

**Week 1-2: Foundation**
- Day 1: Project setup, database schema
- Day 2: Basic SDK structure
- Day 3: Event logging
- Day 4: API endpoints
- Day 5: Testing and documentation
- Day 6-7: Buffer for issues

**Week 3-4: Core Features**
- Day 8: Risk engine heuristics
- Day 9: Policy engine
- Day 10: SDK integrations (LangChain)
- Day 11: Dashboard setup
- Day 12: Dashboard components
- Day 13-14: Integration testing

**Week 5-6: Polish and Launch**
- Day 15: Documentation
- Day 16: Examples and tutorials
- Day 17: Deployment setup
- Day 18: Beta testing
- Day 19: Bug fixes
- Day 20: Launch

### Focus Areas

**High Priority** (Must have for V1):
- SDK wrapper functionality
- Event logging
- Basic risk scoring
- Simple policy enforcement
- Minimal dashboard

**Medium Priority** (Nice to have):
- Advanced risk models
- Complex policies
- Analytics dashboard
- API documentation

**Low Priority** (Defer to V2):
- Multi-tenancy
- Advanced analytics
- Third-party integrations
- Mobile app

### Success Metrics

**V1 Launch Goals**:
- 10 beta users
- 5 agents monitored
- 10,000 events logged
- <5 critical bugs
- Documentation complete

**3-Month Goals**:
- 50 active users
- 100 agents monitored
- 1M events logged
- 10 paying customers
- $1,000 MRR

**6-Month Goals**:
- 200 active users
- 500 agents monitored
- 10M events logged
- 50 paying customers
- $10,000 MRR

---

## Technology Stack Decisions

### Backend Stack

| Component | Technology | Reasoning |
|-----------|-----------|-----------|
| API Framework | FastAPI | Modern, async, auto-docs, type hints |
| Language | Python 3.11+ | AI ecosystem, async support, developer productivity |
| Database | PostgreSQL 15+ | ACID compliance, JSONB, mature ecosystem |
| ORM | SQLAlchemy 2.0 | Async support, type hints, migrations |
| Migrations | Alembic | Standard for SQLAlchemy |
| Caching | Redis 7+ | In-memory speed, pub/sub, simple |
| Queue | Redis (V1), Kafka (V2+) | Start simple, scale later |
| Logging | Loguru | Structured, easy to use |
| Metrics | Prometheus | Industry standard, Grafana integration |
| Testing | Pytest | Async support, fixtures, plugins |

### Frontend Stack

| Component | Technology | Reasoning |
|-----------|-----------|-----------|
| Framework | React 18 | Component model, ecosystem, hiring |
| Build Tool | Vite | Fast, modern, simple config |
| Language | TypeScript | Type safety, better DX |
| Styling | Tailwind CSS | Utility-first, fast development |
| State | TanStack Query | Server state management, caching |
| Routing | React Router | Standard, well-documented |
| Charts | Recharts | React-native, composable |
| Forms | React Hook Form | Performance, validation |

### Infrastructure Stack (V1)

| Component | Technology | Reasoning |
|-----------|-----------|-----------|
| Hosting | Docker Compose | Simple, reproducible, local dev |
| Database | PostgreSQL (Docker) | Easy setup, no cloud costs |
| Reverse Proxy | Nginx | Standard, reliable |
| SSL | Let's Encrypt | Free, automated |
| Monitoring | Prometheus + Grafana | Open source, powerful |

### Infrastructure Stack (V2+)

| Component | Technology | Reasoning |
|-----------|-----------|-----------|
| Cloud | AWS | Mature, comprehensive services |
| Compute | ECS Fargate | Serverless containers, auto-scaling |
| Database | RDS PostgreSQL | Managed, backups, HA |
| Cache | ElastiCache Redis | Managed, reliable |
| Queue | MSK (Kafka) | Managed Kafka, scalable |
| Storage | S3 | Object storage for archives |
| CDN | CloudFront | Global distribution |
| Monitoring | CloudWatch + Datadog | Comprehensive observability |

---

## Conclusion

This documentation provides a comprehensive blueprint for building the AI Control Layer from concept to production. The phased approach ensures we can deliver value quickly while building toward a scalable, enterprise-grade platform.

### Key Takeaways

1. **Start Simple**: V1 focuses on core functionality with minimal complexity
2. **Design for Scale**: Architecture supports future growth without rewrites
3. **Iterate Rapidly**: Ship small increments, gather feedback, improve
4. **Build Moats**: Focus on data, integrations, and compliance for defensibility
5. **Execute Efficiently**: Optimize for limited founder time with clear priorities

### Next Steps

1. Set up development environment
2. Initialize project structure
3. Implement database schema
4. Build SDK core functionality
5. Create API endpoints
6. Develop basic dashboard
7. Write documentation
8. Launch beta program

### Resources

- **GitHub Repository**: [https://github.com/kevinkiplangat432/agents-control-infra-start-up-level-infra.git]
- **Documentation Site**: []
- **API Reference**: []
- **Community Forum**: []
- **Support Email**: kiplangatkevin335@gmail.com

---

**Document Version**: 1.0  
**Last Updated**: 2024-01-15  
**Maintained By**: Kevin email:kiplangatkevin335@gmail.com
**Review Cycle**: Monthly

