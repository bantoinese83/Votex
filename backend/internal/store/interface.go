package store

import "time"

// StoreInterface defines the interface for data access operations
type StoreInterface interface {
	// User operations
	CreateUser(id, username, email, passwordHash string) error
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(id string, updates map[string]interface{}) error
	DeleteUser(id string) error

	// Session operations
	CreateSession(id, userID string, expiresAt string) error
	GetSession(id string) (*Session, error)
	DeleteSession(id string) error

	// Password reset operations
	CreatePasswordResetToken(id, userID, token string, expiresAt time.Time) error
	GetPasswordResetToken(token string) (*PasswordResetToken, error)
	MarkPasswordResetTokenUsed(id string) error
	CleanupExpiredPasswordResetTokens() error
}
