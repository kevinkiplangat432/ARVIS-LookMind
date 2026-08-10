<!-- markdownlint-disable MD033 MD041 MD036-->
<p align="center">
  <img src="../LookMind.png" alt="LookMind" width="300" />
</p>

# Project History

> This is the origin story of ARVIS — why it exists, how it started, and why it changed.
> For where the project is today, see [`README.md`](../README.md).

---

## The Original Question

ARVIS did not begin as a proxy.

It started on **23 February 2026** from a single observation:

As AI agents become more autonomous, organisations gain productivity but lose visibility. Agents make decisions, call tools, access data, and interact with external systems — yet most organisations have no reliable way to monitor, audit, or govern what those agents are actually doing.

The question that started everything:

> How can organisations safely deploy autonomous AI systems while maintaining visibility, accountability, and control?

---

## Phase 1 — AI Control Layer SDK

**February 2026 – April 2026**

The first implementation was an SDK.

The project was built as the **AI Control Layer** — a Python middleware library designed to be embedded directly inside AI agents and agent frameworks. It would intercept agent behaviour and expose observability hooks that dashboards, analytics pipelines, and compliance tooling could consume.

The SDK was built to provide:

- Structured logging of agent actions
- Event streaming and telemetry
- Risk evaluation and scoring
- Policy enforcement
- Kill-switch mechanisms
- Audit trails for compliance

It explored deep questions about AI governance:

- What constitutes risky AI behaviour?
- How should policies be enforced without breaking agent workflows?
- What information is legally required for an audit trail?
- How do you govern an agent you do not control?

The conceptual foundations established in this phase are still present in ARVIS today: observability, governance, auditability, risk detection, append-only audit logs.

---

## The Pivot

A fundamental limitation became clear during development.

**An SDK requires adoption.**

Every application, every agent, every framework, every team must explicitly integrate it before governance can occur. This creates friction at the point of maximum resistance — engineering teams already building with AI tooling do not want to add another dependency before they can ship.

The insight:

> The best governance layer is one that organisations do not have to manually integrate.

Instead of embedding control inside every agent, control could sit at the infrastructure boundary — between the application and the AI provider. Every request passes through regardless of what SDK the agent uses, what language it is written in, or how much the engineering team cares about governance.

This changed everything about the architecture.

---

## Phase 2 — ARVIS

**April 2026 – Present**

The project became **ARVIS — AI Request Visibility & Intelligence System**.

The SDK was abandoned. The Python codebase was set aside. The new system is a Go reverse proxy.

Every request to any LLM API passes through ARVIS. The proxy intercepts it, records it, analyses it, and surfaces the result in a dashboard — with zero changes required from the application sending the request.

The architectural shift delivered several properties the SDK could never offer:

- **Zero adoption friction** — redirect a base URL, done
- **Language-agnostic** — governs Python agents, Node.js agents, Go agents, anything
- **Centralized audit trail** — one Postgres database, one source of truth
- **Vendor-agnostic** — works in front of OpenAI, Anthropic, Mistral, any OpenAI-compatible endpoint

The original goals did not change. The implementation strategy did.

---

## What Exists Today

The current system is documented in full in [`DOCUMENTATION.md`](../DOCUMENTATION.md). In brief:

- A Go binary running two HTTP servers: a proxy on `:8080` and a management API on `:8081`
- PostgreSQL storing every intercepted request and every triggered anomaly
- Three anomaly detection rules: high latency, high token usage, upstream 5xx errors
- A React dashboard polling the API every 5 seconds

The immediate roadmap — PII detection, a Redis-backed policy engine, per-client identity, a kill switch — is described in [`docs/vision.md`](vision.md).

---

## Philosophy

The history of ARVIS reflects a principle that applies to any infrastructure problem:

> Preserve the problem. Evolve the solution.

The problem has been the same since 23 February 2026:

Organisations deploying AI systems need visibility, auditability, and control — at the infrastructure level, without depending on the cooperation of every team and every codebase they operate.

The architecture moved from SDK to proxy. The mission did not move at all.

<!-- markdownlint-disable MD041 -->