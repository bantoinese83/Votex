-- Drop trigger and function
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_password_reset_token_user_id;
DROP INDEX IF EXISTS idx_password_reset_token_token;
DROP INDEX IF EXISTS idx_password_reset_token_expires_at;
DROP INDEX IF EXISTS idx_password_reset_token_used;

-- Drop password reset tokens table
DROP TABLE IF EXISTS password_reset_token;

-- Remove columns from users table
ALTER TABLE "user" DROP COLUMN IF EXISTS email;
ALTER TABLE "user" DROP COLUMN IF EXISTS created_at;
ALTER TABLE "user" DROP COLUMN IF EXISTS updated_at; 