package store

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Age          *int   `db:"age"`
}

type Store struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateUser(id, username, passwordHash string) error {
	query := `INSERT INTO "user" (id, username, password_hash) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, id, username, passwordHash)
	return err
}

func (s *Store) GetUserByID(id string) (*User, error) {
	var user User
	query := `SELECT id, username, password_hash, age FROM "user" WHERE id = $1`
	err := s.DB.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, username, password_hash, age FROM "user" WHERE username = $1`
	err := s.DB.Get(&user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CreateSession(id, userID string, expiresAt string) error {
	query := `INSERT INTO session (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, id, userID, expiresAt)
	return err
}

func (s *Store) GetSession(id string) (*Session, error) {
	var session Session
	query := `SELECT id, user_id, expires_at FROM session WHERE id = $1`
	err := s.DB.Get(&session, query, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) DeleteSession(id string) error {
	query := `DELETE FROM session WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

type Session struct {
	ID        string `db:"id"`
	UserID    string `db:"user_id"`
	ExpiresAt string `db:"expires_at"`
}
