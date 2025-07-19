-- Add email column to users table
ALTER TABLE "user" ADD COLUMN email TEXT UNIQUE;
ALTER TABLE "user" ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE "user" ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP;

-- Create password reset tokens table
CREATE TABLE IF NOT EXISTS password_reset_token (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_email ON "user" (email);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_user_id ON password_reset_token (user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_token ON password_reset_token (token);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_expires_at ON password_reset_token (expires_at);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_used ON password_reset_token (used);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_user_updated_at
    AFTER UPDATE ON "user"
    FOR EACH ROW
    BEGIN
        UPDATE "user" SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END; 