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
1. **Document the "WHY" and "Trade-offs" Never just the what**  reason is anyone can just read the code to see what was built. The **docs/** folder exists to record why it was built that way and what alternative options were rejected.

2. **Provide a 30-second Mental Model (Visula + Text)** Every docs/ feature should allow a new hire to grasp the whole system in under a minute before reading deep details.
include: 
- The core job, what exact business problem does this subsystem solve?
- A System Flow Diagram  showing how data enters, gets processed and where it lands.
- Primary Owners & dependencies: what does it touch (stripe, redis etc)