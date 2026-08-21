<!-- markdownlint-disable-->
<p align="center">
  <img src="LookMind.png" alt="LookMind" width="300" />
</p>

# Automated Runtime Visibility & Intelligence System(ARVIS)

> **Version**: 0.7.0 (Pre-release) · **Status**: Active Development
> **Compliance target**: Kenya Data Protection Act 2019 · Data residency enforced at the infrastructure level

Building infrastructure that lets African financial institutions use
AI without giving up visibility, control, or the ability to prove any
of it happened.

---

## The Problem

Banks, telcos, and fintechs across Kenya are adopting AI faster than
they can govern it. Employees paste customer data into public models.
Nobody can say, with evidence, what left the building and where it
went. And when a regulator asks for proof of compliance, the honest
answer at most organizations is "we don't fully know."

ARVIS exists to close that gap, at the infrastructure level, without
requiring every team to change how they work.

## What It Does

ARVIS sits between your organization and the AI providers you already
use. Point your existing tools at ARVIS instead of directly at the
provider, and every request is intercepted, recorded, and checked
against policy before it goes anywhere. No SDK to install, no code to
rewrite, no workflow to relearn.

Three things fall out of that single design choice, and they map to
three words that get used loosely elsewhere but mean something
specific here:

- **Governance** — the rules of engagement, decided in advance. What's
  allowed, what's blocked, what triggers an immediate stop. This is
  ARVIS acting, before something goes wrong.
- **Auditing** — the record, kept without exception. Every request,
  every flag, written once and never altered. This is ARVIS
  remembering, so nothing has to be reconstructed from memory later.
- **Compliance** — governance and auditing, pointed at a specific
  legal standard and demonstrated against it. This is ARVIS proving
  it, in a form a regulator or an auditor can actually hold.

All three sit on top of a fourth, quieter layer: the raw ability to
see traffic at all. Nothing above it works without it.

## Deployment

ARVIS is built to run entirely inside your own infrastructure, not
ours. A single binary, a local database, no dependency on external
services to function. For an institution that legally cannot put
audit data anywhere it doesn't directly control, that's not a nice
extra, it's the whole reason a tool like this is usable at all.

See [`docs/deployment.md`](docs/deployment.md) for environment setup
across local development, pilot, and production.

## Where This Is Headed

ARVIS today watches individual requests. The direction it's growing
toward is watching *relationships* between requests, which matters a
lot once AI agents start calling other agents, tools, and data
sources on their own, not just answering a single prompt.

The research effort behind this, detailed separately in the project's
research documentation, is exploring three ideas together:

**Understanding structure, not just events.** Instead of treating
every request as an isolated event, the system is being taught to
represent an organization's AI activity as a web of relationships,
who called what, what triggered what next. A technique from the graph
learning field lets a model reason over that web directly, which
matters because a lot of real risk (a chain of agent actions that
individually look fine but add up to something concerning) is only
visible in the relationships, never in any single request on its own.

**Knowing what it doesn't know.** Any model can produce a confident-
looking score even on something it's never actually seen before.
That's a real problem for a compliance tool specifically, a
confidently wrong "this is fine" is worse than an honest "this needs
a human to look at it." So alongside the structural modeling, the
research includes a way for the system to express genuine uncertainty,
not just a risk number, but a sense of how much to trust that number.
High confidence, act on it. Low confidence, hand it to a person
instead of guessing.

**Learning where to look, not what to do.** As the volume of
agent-to-agent activity grows, there's no realistic way to scrutinize
everything with equal depth in real time. The direction being explored
here uses a learning technique to decide where the system's limited
attention should go next, which part of the activity looks most worth
a closer look right now, given what's uncertain and what's changed.
Deliberately, this technique is scoped to deciding *where to pay
attention*, never to deciding *what action to take*. That boundary is
by design, not by policy that could be relaxed later: the layer that
can act (governance) and the layer that decides what's worth watching
more closely stay structurally separate, so ARVIS's ability to observe
faithfully never becomes entangled with a decision it made itself.

None of this changes what ARVIS commits to being: a system that
watches and proves, not one that manages or acts on your behalf.
Growing smarter about *where to look* is not the same as growing
authority to *decide and act*, and that distinction is treated as a
hard architectural line, not a soft guideline.

## Why This, Why Here

Every serious AI governance platform available today was built for a
different regulatory world, the EU AI Act, US frameworks, generic
global compliance. None were built around the Kenya Data Protection
Act 2019, or the operating reality of Kenyan financial institutions.
ARVIS is.

That's a narrow starting point on purpose. The plan isn't to compete
with global platforms on breadth from day one, it's to become the
trusted, evidently-working answer inside this specific regulatory
environment first, and expand from a position of proven credibility
rather than promised scope.

## Status

1.0.0 means real and working, deliberately short of the full
promise. Core interception, logging, and anomaly detection are
functional. Per-identity attribution, the policy engine, PII
detection, and the research direction above are active work, not
finished features. This project would rather be honest about that
than round up to 1.0 early.

`

## Further Reading

| Document | Purpose |
|---|---|
| [`docs/HISTORY.md`](docs/HISTORY.md) | Origin story — why this exists and how it evolved |
| [`docs/vision.md`](docs/vision.md) | The long-term vision for what ARVIS is being built toward |
| [`DOCUMENTATION.md`](DOCUMENTATION.md) | Full technical reference |

## License

This project is proprietary and confidential. All rights reserved.
No part of this repository may be reproduced, distributed, or used to
create derivative works without prior written permission from the
author.

See [`LICENSE`](LICENSE) for the full text.

## Contact

Kevin — <kiplangatkevin335@gmail.com>

Building infrastructure that lets African financial institutions use
AI without giving up visibility, control, or the ability to prove any
of it happened.