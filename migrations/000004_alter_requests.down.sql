ALTER TABLE requests
    DROP COLUMN IF EXISTS identity_id,
    DROP COLUMN IF EXISTS provider;