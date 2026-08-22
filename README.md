<!--markdownlint-disable-->
<p align="center">
  <img src="LookMind.png" alt="LookMind" width="300" />
</p>

# Automated Runtime Visibility & Intelligence System(ARVIS)

**Version 0.9.0 · Active Development · Compliance target: Kenya Data Protection Act 2019**

Building infrastructure that lets African financial institutions use AI without giving up visibility, control, or the ability to prove any of it happened.

---

## The Questions Behind This

Before there was a proxy, a schema, or a line of Go, there were questions that don't have good answers at most organizations right now:

Which employees are sending customer data to a public AI model, and has it happened before? If a regulator asked for proof of compliant AI usage tomorrow, could that proof be produced in an hour, or would it take three weeks? If someone gets fired today, is their access to every AI tool they used actually gone, or just gone from the systems anyone remembered? When something does go wrong, is the first move an honest, evidence-backed answer, or a scramble to reconstruct what even happened?

ARVIS exists because, for almost every organization adopting AI right now, the honest answer to those questions is "we don't fully know." That's the problem. Everything below is the answer.

## What ARVIS Does

ARVIS sits between an organization and the AI providers it already uses. Existing tools point at ARVIS instead of pointing directly at a provider, and from that moment, every request is seen, recorded, and checked before it goes anywhere. No SDK to install. No workflow to relearn. The visibility exists at the infrastructure level, not because every team agreed to adopt a new tool.

## Governance, Auditing, and Compliance, Defined Precisely

These three words get used loosely elsewhere. Here, they mean three specific, separable things, sitting on one shared foundation.

**Observability** is the foundation underneath all three, the raw ability to see AI traffic at all, in real time. Nothing above it functions without it.

**Governance** is control before the fact, the rules of engagement decided in advance: what's allowed, what's blocked, what stops a request mid-flight. This is ARVIS acting, to prevent a problem, not react to one.

**Auditing** is the record kept after the fact, without exception. Every request, every flag, written once and never altered. This is ARVIS remembering, so nothing has to be reconstructed from memory when it matters most.

**Compliance** is governance and auditing pointed at a specific legal standard and demonstrated against it. This is ARVIS proving it, in a form a regulator can actually hold.

## Built to Run Inside Your Own Walls

ARVIS is a single binary and a local database, with nothing built in that depends on an external service to function. For an institution that cannot legally place audit data anywhere outside its own control, that isn't a configuration option, it's the difference between a tool that's usable and one that isn't. Where most platforms in this space treat on-premises deployment as a retrofit onto a cloud-first product, it's the native shape of this one.

## Where the Intelligence Is Headed

ARVIS today watches individual requests. The research direction underneath it, developed in a companion project, is building toward something that watches *relationships between* requests, which matters increasingly as AI stops being one prompt and one answer, and starts being agents calling other agents, tools, and data sources on their own.

Three ideas, being developed together, not separately:

**Seeing the shape of activity, not just isolated events.** Instead of judging each request alone, the system is being taught to represent an organization's AI activity as a web of relationships, who called what, what happened next. Real risk often only shows up in that shape, a chain of individually unremarkable actions that add up to something worth stopping, never visible in any single event on its own.

**Knowing what it doesn't know.** Any model can produce a confident-sounding score on something it has never actually seen before. That's dangerous in a compliance tool specifically, a confidently wrong "this is fine" is worse than an honest "this needs a person to look." So the system is being built to express genuine uncertainty alongside every judgment, not just a risk number, but a sense of how much that number should be trusted. High confidence, act. Low confidence, hand it to a human instead of guessing.

**Learning where attention matters most.** As activity scales, nothing can be scrutinized with equal depth in real time. This piece decides where the system's limited attention goes next, informed by what's uncertain and what's changed. By design, this piece decides only *where to look*, never *what to do about it*. That boundary is architectural, not a policy that could quietly loosen later, the part of the system that can act stays permanently separate from the part that decides what deserves a closer look.

Full detail on this research lives in its own repository: [Adaptive Multi-Agent Reinforcement Learning](https://github.com/kevinkiplangat432/Adaptive-Multi-Agent-Reinforcement-learning).

## Why This, Why Here

Every established AI governance platform available today was built for a different regulatory world, the EU AI Act, US-centric frameworks, generic global compliance retrofitted onto a cloud product. None were built around the Kenya Data Protection Act 2019, or the operating reality of Kenyan financial institutions. This one is.

That's narrow on purpose. The goal isn't to compete on breadth against global platforms from day one, it's to become the proven, trusted answer inside this specific regulatory environment first, and grow outward from that credibility rather than from promised scope.

## Status

0.8.0 reflects what's real and working today, not what's promised. The next milestone is 0.90.0. Release, 1.0.0, comes when the full system works end to end, not on a calendar date chosen in advance.

## Further Reading

| Resource | Purpose |
|---|---|
| [`docs/HISTORY.md`](docs/HISTORY.md) | Origin story, why this exists and how it evolved |
| [`docs/vision.md`](docs/vision.md) | The long-term vision for what ARVIS is being built toward |
| [`DOCUMENTATION.md`](DOCUMENTATION.md) | Full technical reference |
| [Adaptive Multi-Agent RL](https://github.com/kevinkiplangat432/Adaptive-Multi-Agent-Reinforcement-learning) | The GNN, Bayesian uncertainty, and reinforcement learning research behind ARVIS's future intelligence layer |

## License

This project is proprietary and confidential. All rights reserved. No part of this repository may be reproduced, distributed, or used to create derivative works without prior written permission from the author.

See [`LICENSE`](LICENSE) for the full text.

## Contact

Kevin — <kiplangatkevin335@gmail.com>

Building infrastructure that lets African financial institutions use AI without giving up visibility, control, or the ability to prove any of it happened.

