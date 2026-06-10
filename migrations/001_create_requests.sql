CREATE TABLE IF NOT EXISTS requests (
    id                  TEXT PRIMARY KEY,
    model               TEXT,
    prompt_tokens       INTEGER DEFAULT 0,
    completion_tokens   INTEGER DEFAULT 0,
    latency_ms          INTEGER DEFAULT 0,
    status_code         INTEGER DEFAULT 200,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_requests_created ON requests(created_at DESC);
