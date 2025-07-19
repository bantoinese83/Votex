-- Drop trigger
DROP TRIGGER IF EXISTS update_user_updated_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_password_reset_token_user_id;
DROP INDEX IF EXISTS idx_password_reset_token_token;
DROP INDEX IF EXISTS idx_password_reset_token_expires_at;
DROP INDEX IF EXISTS idx_password_reset_token_used;

-- Drop password reset tokens table
DROP TABLE IF EXISTS password_reset_token;

-- Note: SQLite doesn't support DROP COLUMN in older versions
-- For SQLite, we would need to recreate the table without these columns
-- This is a simplified down migration 