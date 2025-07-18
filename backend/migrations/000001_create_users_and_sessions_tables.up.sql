CREATE TABLE IF NOT EXISTS "user" (
    id TEXT PRIMARY KEY,
    age INTEGER,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS session (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES "user"(id),
    expires_at TIMESTAMPTZ NOT NULL
);
