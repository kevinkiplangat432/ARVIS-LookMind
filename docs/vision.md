<!-- markdownlint-disable MD033 MD041 -->

<p align="center">
  <img src="../LookMind.png" alt="LookMind" width="120" />
</p>

# ARVIS — Vision

> *What we are building, why it matters, and where it is going.*

---

## The Problem

Organisations are deploying AI — LLM-powered assistants, agents, internal tools, automated workflows — faster than they can govern it.

The gap is not in the models. The gap is in the infrastructure around the models.

Right now, when a company sends a prompt to an LLM API, they have no reliable answer to any of these questions:

- What data was in that request?
- Did it contain a customer's national ID, phone number, or financial reference?
- Who sent it — which employee, which service, which agent?
- How much did it cost?
- Did the model behave unusually?
- Can we prove, to a regulator, what was sent and when?

The data is leaving the organisation. There is no record of it. There is no gate in front of it. There is nothing watching it.

This is not a hypothetical risk. It is happening today, in every organisation that has deployed AI without a governance layer.

---

## What ARVIS Is

ARVIS is the governance layer that sits between an organisation's systems and any LLM API.

Every request passes through it. Nothing reaches the model without being seen first.

At its core, ARVIS does three things:

**1. Observe** — every request is intercepted and logged. Timestamp, model, token count, latency, status. A complete, append-only audit trail that cannot be altered after the fact.

**2. Detect** — rules run against every request in real time. Anomalous behaviour — unusual latency, excessive token usage, upstream failures, PII in the payload — is flagged immediately, not discovered during a quarterly review.

**3. Enforce** — policy is applied before the request is forwarded. Token budgets, blocked topics, PII masking, kill switches. If a request violates policy, it is stopped at the infrastructure boundary — not after the data has already left.

---

## Why Infrastructure, Not SDK

The first version of this system was a Python SDK. It was abandoned.

An SDK requires adoption. Every team, every application, every agent framework has to integrate it manually before a single request can be governed. Governance only works where the SDK has been deployed. Everywhere else is a blind spot.

A proxy requires nothing. Redirect one environment variable — `OPENAI_BASE_URL` — and every request from every application in the organisation is governed immediately. No code changes. No dependency on engineering teams. No blind spots.

This is the architectural decision that makes ARVIS a governance product rather than a developer tool.

---

## The Full Vision

The proxy is the starting point. It is not the destination.

ARVIS is being built toward a future where every organisation operating AI at scale has a single control plane for all of it — regardless of which models they use, which teams built the agents, or which vendors they depend on.

That control plane will provide:

### Sovereign data handling

Sensitive data — Kenyan National IDs, M-PESA references, KRA PINs, personal information of any kind — is detected and tokenised before it leaves the organisation. The model receives a placeholder. The original value stays inside the infrastructure boundary, encrypted at rest. The organisation never loses custody of its data.

### Identity-aware governance

Every request is attributed to a specific identity — an employee, a service account, an autonomous agent. Token spend, anomaly history, and policy violations are all traceable to the source. Governance without identity is surveillance without accountability.

### Real-time policy enforcement

Rules are stored in Redis and evaluated on every request in under a millisecond. Daily token budgets. Blocked topics. Model allowlists. When a company manager updates a rule through the dashboard, it is enforced on the next request — no restart, no deployment, no lag.

### A kill switch

If a response stream violates policy mid-flight — a blocked topic appears, a budget threshold is crossed — the connection is terminated immediately. The partial response is discarded. The event is logged. The organisation retains control even after the model has started responding.

### Compliance by construction

Audit logs are append-only. Nothing is ever updated or deleted. Every PII detection event, every policy violation, every kill switch trigger is recorded with a timestamp and a reason. A regulator asking for evidence of data governance receives a complete, unalterable record — not a reconstructed narrative.

---

## Who This Is For

**Organisations in regulated industries** — financial services, healthcare, legal, government — that are deploying AI and need to demonstrate compliance with data protection obligations. In Kenya, that means the Data Protection Act 2019. Globally, GDPR, HIPAA, and sector-specific frameworks.

**Enterprises running multiple AI applications** — where different teams use different models, different frameworks, different vendors. ARVIS is vendor-agnostic. It governs OpenAI, Anthropic, Mistral, or any OpenAI-compatible endpoint through the same proxy.

**Engineering and security teams** who need visibility into what their AI systems are doing in production — not as a post-mortem exercise, but in real time, with the ability to act.

**Compliance and risk officers** who are responsible for AI governance but have no current mechanism to enforce or evidence it.

---

## Guiding Principles

**Observability first.** If it cannot be observed, it cannot be trusted. Every architectural decision starts from the question: does this make the system more or less visible?

**Infrastructure over adoption.** Governance that requires manual integration will always have gaps. Governance at the infrastructure boundary has none.

**Real-time over retrospective.** A log you read after the incident happened is an audit trail. A system that detects and stops the incident as it happens is governance. ARVIS is built to be the latter.

**Append-only truth.** The audit trail is the product. Once written, it is never changed. An organisation's ability to prove what happened depends entirely on the integrity of the log.

**Fail open, log everything.** If ARVIS cannot write to the database, it continues proxying. The application never breaks because of governance infrastructure. But the failure is logged, and the gap is visible.

**No lock-in.** ARVIS does not replace the LLM. It sits in front of it. If an organisation stops using ARVIS, they change one environment variable and everything works as before. Governance should be a layer organisations choose to keep because it is valuable — not because removing it is difficult.

---

## The Destination

ARVIS is being built to become the standard infrastructure layer for AI governance in the markets it serves — starting with East Africa, targeting the broader emerging-market enterprise stack where AI adoption is accelerating but governance tooling does not yet exist.

The long-term product is a control plane that:

- Governs AI traffic across an entire organisation from a single deployment
- Provides compliance evidence automatically, without manual report generation
- Detects behavioural anomalies before they become incidents
- Enforces data sovereignty at the infrastructure level
- Scales from a single VPS serving one pilot client to a multi-tenant platform serving enterprises

The immediate product is a proxy that intercepts, logs, and analyses every LLM request — running today, provably working, and ready to put in front of a pilot client.

Both of those things are true at the same time. That is intentional.

---

> *ARVIS is not a monitoring tool bolted onto AI. It is the layer that makes AI safe to deploy at scale.*

---

*See [`README.md`](../README.md) for where the project stands today.*
*See [`docs/HISTORY.md`](HISTORY.md) for how we got here.*
*See [`DOCUMENTATION.md`](../DOCUMENTATION.md) for the technical reference.*
