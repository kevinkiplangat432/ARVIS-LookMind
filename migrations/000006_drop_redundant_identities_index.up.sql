-- key_hash is UNIQUE, which already creates an index. This one was redundant.
DROP INDEX IF EXISTS idx_identities_key_hash; 