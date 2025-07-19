-- Drop indexes
DROP INDEX IF EXISTS idx_session_expires_at;
DROP INDEX IF EXISTS idx_session_user_id;
DROP INDEX IF EXISTS idx_user_username;

-- Drop tables
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS "user"; 