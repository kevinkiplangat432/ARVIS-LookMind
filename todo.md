<!--markdownlint-disable-->
# Auth and Identity: Closing List

## 0. Decide scope (10 minutes, before any code)
These three decisions change the code and the docs, so make them first.

- [ ] Revocation: build it or defer it. Recommendation: build a minimal version
      (revoked_at column, lookup excludes revoked rows, `identity revoke <id>`).
      A leaked key with no kill switch is the first question a bank security
      team asks, and deletion loses attribution (ON DELETE SET NULL).
- [ ] Rate limit on failed auth: defer, record as a known limitation with a
      trigger (before any public exposure or first external customer).
- [ ] Duplicate identity names: allow or reject. Only key_hash is unique today.
      Record the decision either way.

## 1. Documentation standard (unblocks the rest)
- [ ] Add the file header rule to the guide: in a multi file package, a blank
      line separates the file header from `package x`.
- [ ] Add a lighter header form for test files: purpose, which rules it
      guards, how to run it, required environment (for example
      TEST_DATABASE_URL). Decide now so every test file follows one shape.

## 2. keys package
- [ ] EqualHashes (review issue 1): hash both inputs with SHA-256, then
      ConstantTimeCompare the fixed 32 byte results. Same boolean behavior,
      no length shortcut. Fix the comment and package rule 4 to match.
      Timing itself cannot be unit tested, so say so in the test header.
- [ ] keys_test.go header (review issue 2).

## 3. proxy package
- [ ] Extract `extractBearer(header) (string, error)` as a pure function.
      Split on the first space, compare the scheme with strings.EqualFold.
      This also gives us a DB free seam for tests and fuzzing.
- [ ] auth_test.go: flip the lowercase case to "accepted", remove the
      duplicate "bearer without space" rows, run gofmt, add the header.
      Test extractBearer directly. Keep the nil-pool test for rejections.
- [ ] proxy.go (review issue 4): if the request context was canceled, return
      quietly without the error level log or a 503. Keep 503 for real
      lookup failures and deadline exceeded.

## 4. store and commands
- [ ] store: add DeleteIdentity (used to compensate if key delivery fails),
      plus its test. If revocation was approved, add revoked_at, the SQL
      filter, a RevokeIdentity function, and tests.
- [ ] commands/identity.go (review issue 3, ambiguous insert): on insert
      error, look the identity up by its generated ID. If the row exists,
      continue. If not, return the error.
- [ ] commands/identity.go (review issue 1, delivery): write through
      cmd.OutOrStdout(), check every write error. On failure, delete the
      identity and return an error. If the delete also fails, print the ID
      to stderr so an operator can remove it by hand.
- [ ] Headers on store and commands test files. Add a test for the delivery
      failure path using a writer that always errors.
- [ ] If revocation was approved, add the `identity revoke` command and test.

## 5. Schema and data hygiene
- [ ] New migration (never edit an applied one): drop idx_identities_key_hash,
      since UNIQUE already provides the index.
- [ ] One time check for whitespace names created before trimming existed:
      SELECT id, name FROM identities WHERE name <> btrim(name) OR btrim(name) = '';
- [ ] If revocation was approved: migration for revoked_at.

## 6. Streaming path (auth surface only)
- [ ] Read streaming.go for exactly one question: does the outbound request
      replace Authorization with the provider key, as the non streaming path
      does? The "ARVIS key never leaves ARVIS" invariant must hold on both.
- [ ] Write its general package header, a coverage line naming auth, and the
      comment on the auth relevant lines only.

## 7. Verify
- [ ] gofmt, go vet, go test -race ./...
- [ ] Run the database tests with TEST_DATABASE_URL set. Confirm the
      duplicate key_hash test passes against the real schema.
- [ ] Add and run the fuzz targets (see below) for about a minute each.

## 8. Wrap up docs
- [ ] docs/auth_and_identity.md, final, with these facts:
      key_hash UNIQUE, requests keep history but lose attribution on identity
      delete, DB outage returns 503 (review issue 5), lowercase bearer is
      accepted, canceled lookups are not outages, terminal scrollback warning
      for the one time key.
- [ ] Monitoring section: 401 spikes mean probing or a bad client, a rise in
      503 plus the "authentication lookup failed" log means database trouble.
- [ ] Known limitations and deferred items section (below).
- [ ] docs/README.md coverage index row: auth and identity, done.
- [ ] Push, let the reviewer run again, confirm no new blocking findings.

## Deferred on purpose (recorded in the doc, not forgotten)
- Rate limiting failed auth (trigger above)
- store.ErrNotFound so the proxy stops depending on pgx directly
- A lookup interface so authenticate's database paths are testable without
  Postgres
- Key rotation tooling beyond "create a new identity"
- Human SSO (OIDC and SAML)

## Done means
Every box above is checked, or sits under Deferred with a trigger. Nothing
in the docs contradicts the code, and every package rule has a test behind it.