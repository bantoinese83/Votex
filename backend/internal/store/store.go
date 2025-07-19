package store

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

// Error definitions
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrTokenNotFound = errors.New("password reset token not found")
	ErrTokenExpired  = errors.New("password reset token expired")
)

type User struct {
	ID           string     `db:"id"`
	Username     string     `db:"username"`
	Email        *string    `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Age          *int       `db:"age"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type PasswordResetToken struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Token     string     `db:"token"`
	ExpiresAt time.Time  `db:"expires_at"`
	Used      bool       `db:"used"`
	CreatedAt *time.Time `db:"created_at"`
}

type Store struct {
	DB       *sqlx.DB
	IsSQLite bool
}

func New(db *sqlx.DB, isSQLite bool) *Store {
	return &Store{DB: db, IsSQLite: isSQLite}
}

func (s *Store) CreateUser(id, username, email, passwordHash string) error {
	var query string
	if s.IsSQLite {
		query = `INSERT INTO "user" (id, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))`
	} else {
		query = `INSERT INTO "user" (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())`
	}
	_, err := s.DB.Exec(query, id, username, email, passwordHash)
	return err
}

func (s *Store) GetUserByID(id string) (*User, error) {
	var user User
	var query string
	if s.IsSQLite {
		query = `SELECT id, username, email, password_hash, age, created_at, updated_at FROM "user" WHERE id = ?`
	} else {
		query = `SELECT id, username, email, password_hash, age, created_at, updated_at FROM "user" WHERE id = $1`
	}
	err := s.DB.Get(&user, query, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	var user User
	var query string
	if s.IsSQLite {
		query = `SELECT id, username, email, password_hash, age, created_at, updated_at FROM "user" WHERE username = ?`
	} else {
		query = `SELECT id, username, email, password_hash, age, created_at, updated_at FROM "user" WHERE username = $1`
	}
	err := s.DB.Get(&user, query, username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	var user User
	var query string
	if s.IsSQLite {
		query = `SELECT id, username, email, password_hash, age, created_at, updated_at FROM "user" WHERE email = ?`
	} else {
		query = `SELECT id, username, email, password_hash, age, created_at, updated_at FROM "user" WHERE email = $1`
	}
	err := s.DB.Get(&user, query, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s *Store) UpdateUser(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// Add updated_at timestamp
	updates["updated_at"] = "NOW()"

	var query string
	var args []interface{}
	argIndex := 1

	// Build dynamic query
	setClause := ""
	for field, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		if s.IsSQLite {
			setClause += field + " = ?"
		} else {
			setClause += field + " = $" + string(rune('0'+argIndex))
		}
		args = append(args, value)
		argIndex++
	}

	if s.IsSQLite {
		query = `UPDATE "user" SET ` + setClause + ` WHERE id = ?`
	} else {
		query = `UPDATE "user" SET ` + setClause + ` WHERE id = $` + string(rune('0'+argIndex))
	}
	args = append(args, id)

	_, err := s.DB.Exec(query, args...)
	return err
}

func (s *Store) DeleteUser(id string) error {
	var query string
	if s.IsSQLite {
		query = `DELETE FROM "user" WHERE id = ?`
	} else {
		query = `DELETE FROM "user" WHERE id = $1`
	}
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *Store) CreateSession(id, userID string, expiresAt string) error {
	var query string
	if s.IsSQLite {
		query = `INSERT INTO session (id, user_id, expires_at) VALUES (?, ?, ?)`
	} else {
		query = `INSERT INTO session (id, user_id, expires_at) VALUES ($1, $2, $3)`
	}
	_, err := s.DB.Exec(query, id, userID, expiresAt)
	return err
}

func (s *Store) GetSession(id string) (*Session, error) {
	var session Session
	var query string
	if s.IsSQLite {
		query = `SELECT id, user_id, expires_at FROM session WHERE id = ?`
	} else {
		query = `SELECT id, user_id, expires_at FROM session WHERE id = $1`
	}
	err := s.DB.Get(&session, query, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) DeleteSession(id string) error {
	var query string
	if s.IsSQLite {
		query = `DELETE FROM session WHERE id = ?`
	} else {
		query = `DELETE FROM session WHERE id = $1`
	}
	_, err := s.DB.Exec(query, id)
	return err
}

// Password reset token methods
func (s *Store) CreatePasswordResetToken(id, userID, token string, expiresAt time.Time) error {
	var query string
	if s.IsSQLite {
		query = `INSERT INTO password_reset_token (id, user_id, token, expires_at, used, created_at) VALUES (?, ?, ?, ?, ?, datetime('now'))`
	} else {
		query = `INSERT INTO password_reset_token (id, user_id, token, expires_at, used, created_at) VALUES ($1, $2, $3, $4, $5, NOW())`
	}
	_, err := s.DB.Exec(query, id, userID, token, expiresAt, false)
	return err
}

func (s *Store) GetPasswordResetToken(token string) (*PasswordResetToken, error) {
	var resetToken PasswordResetToken
	var query string
	if s.IsSQLite {
		query = `SELECT id, user_id, token, expires_at, used, created_at FROM password_reset_token WHERE token = ?`
	} else {
		query = `SELECT id, user_id, token, expires_at, used, created_at FROM password_reset_token WHERE token = $1`
	}
	err := s.DB.Get(&resetToken, query, token)
	if err != nil {
		return nil, ErrTokenNotFound
	}
	return &resetToken, nil
}

func (s *Store) MarkPasswordResetTokenUsed(id string) error {
	var query string
	if s.IsSQLite {
		query = `UPDATE password_reset_token SET used = ? WHERE id = ?`
	} else {
		query = `UPDATE password_reset_token SET used = $1 WHERE id = $2`
	}
	_, err := s.DB.Exec(query, true, id)
	return err
}

func (s *Store) CleanupExpiredPasswordResetTokens() error {
	var query string
	if s.IsSQLite {
		query = `DELETE FROM password_reset_token WHERE expires_at < datetime('now') OR used = ?`
	} else {
		query = `DELETE FROM password_reset_token WHERE expires_at < NOW() OR used = $1`
	}
	_, err := s.DB.Exec(query, true)
	return err
}

type Session struct {
	ID        string `db:"id"`
	UserID    string `db:"user_id"`
	ExpiresAt string `db:"expires_at"`
}

// MockStore is a mock implementation for development when database is not available
type MockStore struct{}

func (m *MockStore) CreateUser(id, username, email, passwordHash string) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) GetUserByID(id string) (*User, error) {
	// Mock implementation - return a mock user
	return &User{
		ID:           id,
		Username:     "mockuser",
		Email:        stringPtr("mock@example.com"),
		PasswordHash: "$2a$10$mockhash",
	}, nil
}

func (m *MockStore) GetUserByUsername(username string) (*User, error) {
	// Mock implementation - return a mock user for any username
	return &User{
		ID:           "mock-id",
		Username:     username,
		Email:        stringPtr("mock@example.com"),
		PasswordHash: "$2a$10$mockhash",
	}, nil
}

func (m *MockStore) GetUserByEmail(email string) (*User, error) {
	// Mock implementation - return a mock user for any email
	return &User{
		ID:           "mock-id",
		Username:     "mockuser",
		Email:        &email,
		PasswordHash: "$2a$10$mockhash",
	}, nil
}

func (m *MockStore) UpdateUser(id string, updates map[string]interface{}) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) DeleteUser(id string) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) CreateSession(id, userID string, expiresAt string) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) GetSession(id string) (*Session, error) {
	// Mock implementation - return nil (no session found)
	return nil, errors.New("session not found")
}

func (m *MockStore) DeleteSession(id string) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) CreatePasswordResetToken(id, userID, token string, expiresAt time.Time) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) GetPasswordResetToken(token string) (*PasswordResetToken, error) {
	// Mock implementation - return nil (no token found)
	return nil, ErrTokenNotFound
}

func (m *MockStore) MarkPasswordResetTokenUsed(id string) error {
	// Mock implementation - always succeeds
	return nil
}

func (m *MockStore) CleanupExpiredPasswordResetTokens() error {
	// Mock implementation - always succeeds
	return nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
