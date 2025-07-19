# SQLite Fallback Implementation

This backend now supports SQLite as a fallback database when PostgreSQL is not available. This is particularly useful for development environments where PostgreSQL might not be running.

## How it works

1. **Primary Database**: PostgreSQL (configured via `DB_URL`)
2. **Fallback Database**: SQLite (configured via `SQLITE_PATH`)
3. **Automatic Fallback**: If PostgreSQL connection fails, the system automatically falls back to SQLite

## Configuration

### Environment Variables

```env
# Database Configuration
DB_URL=postgres://user:password@localhost:5432/vortexdb?sslmode=disable
DB_TYPE=postgres  # or "sqlite" to force SQLite
SQLITE_PATH=./data/votex.db
```

### Default Behavior

- **Development**: Tries PostgreSQL first, falls back to SQLite if unavailable
- **Production**: Requires PostgreSQL to be available (no fallback)

## Database Structure

Both PostgreSQL and SQLite use the same schema:

### Users Table
```sql
CREATE TABLE "user" (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    age INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Sessions Table
```sql
CREATE TABLE session (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);
```

## Migration Files

Migrations are organized by database type:
- `migrations/postgres/` - PostgreSQL-specific migrations
- `migrations/sqlite/` - SQLite-specific migrations

## Features

- ✅ Automatic fallback from PostgreSQL to SQLite
- ✅ Proper SQL syntax handling for both databases
- ✅ Migration support for both databases
- ✅ Foreign key constraints in SQLite
- ✅ Indexes for performance
- ✅ Development-friendly with mock store as last resort

## Testing

The system has been tested with:
- User registration
- User login
- JWT token authentication
- Profile retrieval

## File Structure

```
backend/
├── data/
│   └── votex.db          # SQLite database file
├── migrations/
│   ├── postgres/         # PostgreSQL migrations
│   └── sqlite/           # SQLite migrations
├── internal/
│   ├── config/
│   │   └── config.go     # Database configuration
│   └── store/
│       ├── store.go      # Database operations
│       └── database.go   # Connection handling
└── cmd/server/
    └── main.go           # Server startup with fallback logic
```

## Usage

1. Start the server: `go run ./cmd/server/main.go`
2. The server will automatically detect available databases
3. Check logs for database connection status
4. API endpoints work the same regardless of database type

## Benefits

- **Development**: No need to run PostgreSQL locally
- **Portability**: SQLite database file can be easily shared
- **Simplicity**: Single file database for development
- **Reliability**: Automatic fallback ensures the app always works 