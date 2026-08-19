ALTER TABLE requests
    ADD COLUMN IF NOT EXISTS identity_id TEXT REFERENCES identities(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS provider TEXT NOT NULL DEFAULT 'unknown';

CREATE INDEX IF NOT EXISTS idx_requests_identity_id ON requests(identity_id);
CREATE INDEX IF NOT EXISTS idx_requests_provider ON requests(provider);