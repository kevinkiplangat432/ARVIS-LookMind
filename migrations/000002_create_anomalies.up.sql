CREATE TABLE IF NOT EXISTS anomalies (
    id          TEXT PRIMARY KEY,
    request_id  TEXT NOT NULL REFERENCES requests(id) ON DELETE RESTRICT,
    rule        TEXT NOT NULL,
    detail      TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_anomalies_created ON anomalies(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_anomalies_request_id ON anomalies(request_id);