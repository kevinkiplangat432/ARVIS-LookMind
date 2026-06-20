CREATE TABLE IF NOT EXISTS anomalies (
    id          TEXT PRIMARY KEY,
    request_id  TEXT REFERENCES requests(id),
    rule        TEXT NOT NULL,
    detail      TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_anomalies_created ON anomalies(created_at DESC);
