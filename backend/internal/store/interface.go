package store

// StoreInterface defines the interface for data access operations
type StoreInterface interface {
	CreateUser(id, username, passwordHash string) error
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	CreateSession(id, userID string, expiresAt string) error
	GetSession(id string) (*Session, error)
	DeleteSession(id string) error
}
