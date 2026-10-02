<!--markdownlint-disable-->
# ARVIS Engineering Documentation Standard

This document was written in 2026, before the first hire, when the whole company was two people and a Postgres container. We wrote the standard first because habits set in the first month become culture by the first year, and culture is very hard to refactor.

If you are reading this, you are part of the team that inherited it. Follow it, argue with it, improve it. Do not ignore it.

## 1. Principles

1. **Docs are part of the code.** A change that alters behavior and leaves the docs wrong is an incomplete change. Reviewers reject it the same way they reject a failing test.
2. **Explain why, not what.** Anyone can read the code to see what was built. The code cannot tell you what we rejected, what we feared, or what broke at 2am in the past. Only we can write that down.
3. **Write for a stranger.** The reader has no context, no Slack history, and a deadline. Be clear first, clever second.
4. **Stale docs are bugs.** A wrong comment is worse than no comment, because it lies with confidence. File it and fix it like any other bug.
5. **Be readable, not boring.** Documentation nobody finishes protects nobody. Voice is allowed, and encouraged, within the rules in section 7.

## 2. How We Work: One Feature, End to End

Do not document three features to finish one file. Follow one feature completely, then move on.

1. Trace the feature end to end: route, handler, store, types, errors.
2. Write its `docs/<feature>.md` first.
3. Write the package header for every package the feature touches, general things only.
4. Comment the functions and sections the feature actually uses.
5. Write or update the tests for every file you consider finished.

If a file is shared by several features, document only the part you traced and say so with a coverage line (section 3). A file must never pretend to be fully documented when it is not.

Full documentation of a file is never a single sitting. It accumulates, one traced feature at a time. That is a feature of the process, not a gap in it.

## 3. Level One: In-File Comments

Every file with code has a package header and comments on its non obvious functions. Not just the security files, not just the famous ones. The boring utility file you wrote on a Friday is exactly the one someone will debug on a Monday.

### The package header

It sits above `package x`. It tells the reader why the file exists and what must never break.

1. **Purpose:** what it is for and why it exists.
2. **Rules:** the invariants that must never be broken.
3. **Dependencies:** what it uses, and what it must never import.
4. **Usage:** who calls it, and how.
5. **Thread safety:** pure and stateless, or holding global state. Safe for concurrent goroutines, or not.
6. **Docs pointer:** the `docs/` file that explains the bigger picture.
7. **Coverage line (partial files only):** which feature is documented so far.

```go
// Package auth provides API key generation and hashing for ARVIS machine
// callers. Human SSO is out of scope.
//
// Rules that must never be broken:
//  1. Raw keys are never logged, persisted, or kept after issuance.
//  2. Only the output of HashKey is stored or queried.
//  3. This package never touches the database, HTTP, or config.
//
// Dependencies: standard library only.
// Usage: cmd/identity at issuance, proxy/auth.go on every request.
// Thread safety: stateless, safe for concurrent use.
//
// Docs: docs/auth_and_identity.md
package auth
```

A partial file carries an honest label:

```go
// Coverage: only the auth path is documented so far.
// See docs/auth_and_identity.md. The rest is uncharted territory.
```

### Function and inline comments

1. Short and descriptive. One or two lines, not an essay.
2. Why, not what. The code already says what.
3. Always comment the surprising: fail open versus fail closed, ordering that matters, security reasoning, anything that looks wrong but is deliberate.
4. Never restate the function name. "HashKey hashes a key" has helped nobody in history.

## 4. Level Two: The docs/ Folder

In-file comments explain the code and its invariants. The `docs/` folder explains how systems fit together. File names say what is inside, for example `auth_and_identity.md`.

Every feature doc covers four things:

1. **The why and the trade-offs.** What we built, why, and which alternatives we rejected and why. This is the part nobody can reconstruct later.
2. **A 30 second mental model.** The core business problem, a flow diagram (Mermaid or ASCII), and the dependencies it touches. A new hire should grasp the whole subsystem in under a minute.
3. **Invariants and failure scenarios.** What must never happen, and what happens when Postgres, Redis, or a provider goes down. For every failure, say whether it fails open or closed, and justify it.
4. **Operating procedures.** How to test it, which log or metric tells you it is failing, and step by step runbooks for common tasks.

Template:

```markdown
# [Feature or Subsystem Name]

## 1. High-Level Overview
- **Purpose:** what this does and why it exists.
- **Key Dependencies:** services, databases, external APIs.

## 2. Architecture and Data Flow
[Mermaid or ASCII diagram]

## 3. Key Technical Decisions and Trade-Offs
- Why we built it this way instead of [Alternative X].
- Known limitations.

## 4. Invariants and Security
- Rules that must never break.
- How sensitive data is handled.

## 5. Failure Modes and Operational Guidance
- What fails if Service X goes down.
- How to test, monitor, and troubleshoot.
- Runbooks for common tasks.
```

Every feature doc is listed in `docs/README.md`, so the state of our documentation is visible at a glance:

```markdown
| Feature | Docs file | Package comments | Tests | Status |
|---|---|---|---|---|
| Auth and identity | auth_and_identity.md | done | keys done | in progress |
```

## 5. Tests

A file is not finished being refactored until it has a `_test.go`.

1. Table driven tests where it fits.
2. Test behavior and invariants, not implementation details.
3. Every rule in a package header gets at least one test that fails if the rule breaks. A rule with no test is a wish.
4. Test the failure paths. Bad input, empty input, and the input the original author never imagined.

## 6. Keeping Docs Alive

1. Behavior change and doc change ship in the same pull request.
2. If you read a doc and it is wrong, fix it in your next commit. Do not leave a note for "someone".
3. When you trace a feature and find the docs already cover it, extend them rather than writing a second version.
4. Delete docs that describe things that no longer exist. History lives in git, not in the docs.

## 7. Voice and Humour

We write with personality. Documentation nobody enjoys reading is documentation nobody reads, and unread docs protect no one. The tradition here is blunt, vivid, and a little irreverent, in the spirit of the best engineering writing: say the true thing clearly, and make it stick.

The rules:

1. **State the rule plainly first.** Humour comes after clarity, never instead of it. A reader in a hurry must get the rule even if they miss the joke.
2. **The delete test.** Remove the joke. If the comment is no longer complete and correct, the joke was carrying weight it should not carry.
3. **The joke must teach.** The best humour makes a rule memorable. If it only shows off, cut it.
4. **Punch at code, mistakes, and our past selves.** Never at people, teams, customers, or anyone's background.
5. **Runbooks stay plain.** During an incident, nobody wants wit. Steps are steps.
6. **Short.** A good joke fits in a line. If it needs a paragraph, it is a blog post.
7. **Keep it timeless.** No jokes that depend on a meme, a news cycle, or an inside reference only three people understand. The reader in five years should still get it.

Examples of the tone we want:

```go
// Keys are shown once. We cannot recover them. We are a security product, not a magician.
```

```go
// Tokenization fails closed. Sending a customer's national ID to a third party is not "degraded service", it is a headline.
```

```go
// Redis down during a policy check: fail open. Annoying, recoverable.
// Do not copy this behavior into tokenization. See above, and then see a therapist.
```

Examples of the tone we do not want:

1. A riddle where the invariant should be.
2. Sarcasm aimed at a colleague's old code.
3. A joke so clever that nobody can tell whether the rule is real.

## 8. Reviewer Checklist

Before approving a pull request that touches a feature, confirm:

1. The `docs/` file exists, is accurate, and is listed in the index.
2. Every package touched has a header, and partial files carry a coverage line.
3. Surprising decisions have a short why comment.
4. Files considered finished have tests, and every package rule has a test behind it.
5. Nothing in the docs is now false because of this change.
6. The humour, if present, passes section 7.

## Closing Note

Rules like these are boring to write and priceless to inherit. Somewhere in the future, an engineer will open a file at 2am, find a clear header, a clear why, and a line that makes them smile, and they will fix the problem faster because of it. That engineer is why this document exists.

Write for them. 

**_Author kevin Your CEO or CTO or I Choose a different Path_**