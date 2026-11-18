DROP INDEX IF EXISTS sessions_user_id_idx;
ALTER TABLE sessions DROP COLUMN IF EXISTS revoked_at;
