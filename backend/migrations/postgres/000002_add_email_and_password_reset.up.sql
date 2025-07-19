-- Add email column to users table
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS email TEXT UNIQUE;
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- Create password reset tokens table
CREATE TABLE IF NOT EXISTS password_reset_token (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_email ON "user" (email);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_user_id ON password_reset_token (user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_token ON password_reset_token (token);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_expires_at ON password_reset_token (expires_at);
CREATE INDEX IF NOT EXISTS idx_password_reset_token_used ON password_reset_token (used);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_user_updated_at 
    BEFORE UPDATE ON "user" 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 