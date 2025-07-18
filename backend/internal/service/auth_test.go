package service

import (
	"testing"

	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/store"
)

func TestAuthService_Register(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	// Mock store for testing
	mockStore := &MockStore{}
	service := NewAuthService(mockStore, cfg)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid registration",
			username: "testuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, user, err := service.Register(tt.username, tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Register() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Register() unexpected error: %v", err)
				return
			}

			if token == "" {
				t.Errorf("Register() expected token but got empty string")
			}

			if user == nil {
				t.Errorf("Register() expected user but got nil")
			}

			if user.Username != tt.username {
				t.Errorf("Register() user.Username = %v, want %v", user.Username, tt.username)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	mockStore := &MockStore{}
	service := NewAuthService(mockStore, cfg)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid login",
			username: "testuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid credentials",
			username: "wronguser",
			password: "wrongpass",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, user, err := service.Login(tt.username, tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Login() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Login() unexpected error: %v", err)
				return
			}

			if token == "" {
				t.Errorf("Login() expected token but got empty string")
			}

			if user == nil {
				t.Errorf("Login() expected user but got nil")
			}
		})
	}
}

// MockStore is a mock implementation of the store for testing
type MockStore struct{}

func (m *MockStore) CreateUser(id, username, passwordHash string) error {
	return nil
}

func (m *MockStore) GetUserByID(id string) (*store.User, error) {
	return &store.User{
		ID:           id,
		Username:     "testuser",
		PasswordHash: "$2a$10$hashedpassword",
	}, nil
}

func (m *MockStore) GetUserByUsername(username string) (*store.User, error) {
	if username == "testuser" {
		return &store.User{
			ID:           "123",
			Username:     username,
			PasswordHash: "$2a$10$hashedpassword",
		}, nil
	}
	return nil, nil
}

func (m *MockStore) CreateSession(id, userID string, expiresAt string) error {
	return nil
}

func (m *MockStore) GetSession(id string) (*store.Session, error) {
	return nil, nil
}

func (m *MockStore) DeleteSession(id string) error {
	return nil
}
