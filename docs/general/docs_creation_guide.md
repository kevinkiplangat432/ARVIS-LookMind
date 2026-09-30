<!--markdownlint-disable-->
# A guide to documentations in this codebase.

There are two signature levels for documentation in this codebase unless you have choosen that you should have your own too ( which is okay) but i know two.


## Top-level In-file comments.
This is the second type of documentation.
This are comments that sit at the very beginning of a file/package( **above *package auth* or at the top of a module file**). Their job is to tell a developer reading the code why this file exists and what rules must never be broken.

##### Top level  In-file comments should cover this 4 key areas:
1. **The purpose of the file** - This is a high level overview of what this file is for, what it does and why it exists.
2. **The rules of the file** - This is a list of rules that must never be broken, for instance if this file is a utility file for the auth module then it should not have any references to the database or any other modules.
3. **The dependencies of the file** - This is a list of all the dependencies that 
this file has, for instance if this file is a utility file for the auth module then it should have a dependency on the auth module but not on any other modules.
4. **The usage of the file** - This is a list of all the ways this 
file can be used, for instance if this file is a utility file for the auth module then it should have a usage section that explains how to use the functions in this file.
5. **Thread Safety & Execution Guarantees** Is this code safe for concurrent gouroutines/threads? is it a pure utility package, or does it hold global state?

```go
// Package auth provides cryptographic primitives for API key generation,
// hashing, and request token validation across ARVIS services.
//
// Critical Security Invariants:
//  - Raw API keys MUST NEVER be written to logs, cached in memory beyond initial issuance,
//    or stored in plain text within persistence layers.
//  - Only the deterministic SHA-256 digest (produced by HashKey) may be stored or queried.
//
// Usage & Thread Safety:
//  - All functions in this package are stateless, pure, and safe for concurrent execution.
//  - Shared across both identity management CLI tools and incoming edge proxy middleware.

package auth
```


## The primary documentation  (docs/)
While in-file comments explain code execution and invariants, your **docs/** folder explains how systems fit together.

This is the documentation you are also reading right now it is a file that "talks" about another file and in this codebase we are using the .md file extension. 

What is here is mostly cross-cutting system design,architecture diagrams and high level feature guides. 

The file name often always capture what is talks about, for instance auth_and_identity as you probably have guessed talks about the authentication and identity features.

##### The docs/ should cover this main areas.
1. **Document the "WHY" and "Trade-offs" Never just the what**  reason is, anyone can just read the code to see what was built. The **docs/** folder exists to record why it was built that way and what alternative options were rejected.

2. **Provide a 30-second Mental Model (Visula + Text)** Every docs/ feature should allow a new hire to grasp the whole system in under a minute before reading deep details.
include: 
- The core job, what exact business problem does this subsystem solve?
- A System Flow Diagram  showing how data enters, gets processed and where it lands.
- Primary Owners & dependencies: what does it touch (stripe, redis etc)

3. **Define invariants and failure scenarios** Describe how the system behaves when things go wrong. This how you will save massive amounts of debugging time.
- **Hard invariants** What conditions must never occur under any circumstances? for example "An order can never ne marked "shipped" without a valid Payment ID"

- **Failure modes & Recovery** What happens if the database goes down mid-transaction? What happens if an upstream API times out? Is the Operation retryable or does it dead-letter?

>This examples are obviously not fit to our system so don't go thinking we are amazon.

4. **Specify Operating & Maintenance Procedures** Documentation should provide actionable guidance when a incident occurs or when a team member needs to perform a standard operating task.

Every core feature file in docs/ shouled include:
- **How to test** Specific commands or local setups needed to test this feature manually or via automated suites.
- **Monitoring & alerts** Which metric or error log indicates this feature us failing in production?
- **Common Runbooks/Standard Operating Procedure(SOP)** Step-by-step instructions for common operational task (e.g "How to manually invalidate a customer account" or "How to replay failed webhooks")

The following is a suggested template"
```markdown
# [Feature / Subsystem Name]

## 1. High-Level Overview
- **Purpose:** Brief description of what this does and why it exists.
- **Key Dependencies:** Services, databases, or external APIs involved.

## 2. Architecture & Data Flow
[ Insert Mermaid.js chart or simple ASCII data flow diagram here ]

## 3. Key Technical Decisions & Trade-Offs
- Why we built it this way instead of [Alternative X].
- Known limitations or performance constraints.

## 4. Invariants & Security
- Non-negotiable rules that must never break.
- How sensitive data is handled.

## 5. Failure Modes & Operational Guidance
- What fails if Service X goes down?
- How to test locally and troubleshoot issues.

```

