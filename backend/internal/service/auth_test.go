package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/store"
)

// MockStore is a mock implementation of StoreInterface
type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateUser(id, username, email, passwordHash string) error {
	args := m.Called(id, username, email, passwordHash)
	return args.Error(0)
}

func (m *MockStore) GetUserByUsername(username string) (*store.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) GetUserByEmail(email string) (*store.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) GetUserByID(id string) (*store.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) UpdateUser(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockStore) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreatePasswordResetToken(id, userID, token string, expiresAt time.Time) error {
	args := m.Called(id, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockStore) GetPasswordResetToken(token string) (*store.PasswordResetToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.PasswordResetToken), args.Error(1)
}

func (m *MockStore) MarkPasswordResetTokenUsed(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockStore) UpdateUserPassword(userID, passwordHash string) error {
	args := m.Called(userID, passwordHash)
	return args.Error(0)
}

func (m *MockStore) CleanupExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStore) CleanupExpiredPasswordResetTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStore) CreateSession(id, userID string, expiresAt string) error {
	args := m.Called(id, userID, expiresAt)
	return args.Error(0)
}

func (m *MockStore) GetSession(id string) (*store.Session, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.Session), args.Error(1)
}

func (m *MockStore) DeleteSession(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		email         string
		password      string
		setupMock     func(*MockStore)
		expectedError error
	}{
		{
			name:     "successful registration",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("GetUserByUsername", "testuser").Return(nil, assert.AnError)
				mockStore.On("GetUserByEmail", "test@example.com").Return(nil, assert.AnError)
				mockStore.On("CreateUser", mock.AnythingOfType("string"), "testuser", "test@example.com", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "username already exists",
			username: "existinguser",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockStore *MockStore) {
				existingUser := &store.User{ID: "1", Username: "existinguser"}
				mockStore.On("GetUserByUsername", "existinguser").Return(existingUser, nil)
			},
			expectedError: ErrUserExists,
		},
		{
			name:     "email already exists",
			username: "newuser",
			email:    "existing@example.com",
			password: "password123",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("GetUserByUsername", "newuser").Return(nil, assert.AnError)
				existingUser := &store.User{ID: "1", Email: stringPtr("existing@example.com")}
				mockStore.On("GetUserByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedError: ErrEmailExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			tt.setupMock(mockStore)

			cfg := &config.Config{
				PasswordResetTokenExpiry: 24,
				AppURL:                   "http://localhost:5173",
			}

			service := &AuthService{
				Store:        mockStore,
				Cfg:          cfg,
				EmailService: NewEmailService(cfg),
			}

			token, user, err := service.Register(tt.username, tt.email, tt.password)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, *user.Email)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		setupMock     func(*MockStore)
		expectedError error
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "password123",
			setupMock: func(mockStore *MockStore) {
				hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890"
				user := &store.User{
					ID:           "1",
					Username:     "testuser",
					PasswordHash: hashedPassword,
					Email:        stringPtr("test@example.com"),
				}
				mockStore.On("GetUserByUsername", "testuser").Return(user, nil)
			},
			expectedError: nil,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "password123",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("GetUserByUsername", "nonexistent").Return(nil, assert.AnError)
			},
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			tt.setupMock(mockStore)

			cfg := &config.Config{
				PasswordResetTokenExpiry: 24,
				AppURL:                   "http://localhost:5173",
			}

			service := &AuthService{
				Store:        mockStore,
				Cfg:          cfg,
				EmailService: NewEmailService(cfg),
			}

			token, user, err := service.Login(tt.username, tt.password)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthService_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setupMock     func(*MockStore)
		expectedError error
	}{
		{
			name:  "successful password reset request",
			email: "test@example.com",
			setupMock: func(mockStore *MockStore) {
				user := &store.User{
					ID:       "1",
					Username: "testuser",
					Email:    stringPtr("test@example.com"),
				}
				mockStore.On("GetUserByEmail", "test@example.com").Return(user, nil)
				mockStore.On("CreatePasswordResetToken", mock.AnythingOfType("string"), "1", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "email not found (should not reveal existence)",
			email: "nonexistent@example.com",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("GetUserByEmail", "nonexistent@example.com").Return(nil, assert.AnError)
			},
			expectedError: nil, // Should not return error for security
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			tt.setupMock(mockStore)

			cfg := &config.Config{
				PasswordResetTokenExpiry: 24,
				AppURL:                   "http://localhost:5173",
			}

			service := &AuthService{
				Store:        mockStore,
				Cfg:          cfg,
				EmailService: NewEmailService(cfg),
			}

			err := service.RequestPasswordReset(tt.email)

			assert.Equal(t, tt.expectedError, err)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthService_ResetPassword(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		newPassword   string
		setupMock     func(*MockStore)
		expectedError error
	}{
		{
			name:        "successful password reset",
			token:       "valid-token",
			newPassword: "newpassword123",
			setupMock: func(mockStore *MockStore) {
				resetToken := &store.PasswordResetToken{
					ID:        "1",
					UserID:    "1",
					Token:     "valid-token",
					ExpiresAt: time.Now().Add(time.Hour),
					Used:      false,
				}
				mockStore.On("GetPasswordResetToken", "valid-token").Return(resetToken, nil)
				mockStore.On("MarkPasswordResetTokenUsed", "valid-token").Return(nil)
				mockStore.On("UpdateUserPassword", "1", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "token not found",
			token:       "invalid-token",
			newPassword: "newpassword123",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("GetPasswordResetToken", "invalid-token").Return(nil, assert.AnError)
			},
			expectedError: ErrTokenNotFound,
		},
		{
			name:        "token expired",
			token:       "expired-token",
			newPassword: "newpassword123",
			setupMock: func(mockStore *MockStore) {
				resetToken := &store.PasswordResetToken{
					ID:        "1",
					UserID:    "1",
					Token:     "expired-token",
					ExpiresAt: time.Now().Add(-time.Hour), // Expired
					Used:      false,
				}
				mockStore.On("GetPasswordResetToken", "expired-token").Return(resetToken, nil)
			},
			expectedError: ErrTokenExpired,
		},
		{
			name:        "token already used",
			token:       "used-token",
			newPassword: "newpassword123",
			setupMock: func(mockStore *MockStore) {
				resetToken := &store.PasswordResetToken{
					ID:        "1",
					UserID:    "1",
					Token:     "used-token",
					ExpiresAt: time.Now().Add(time.Hour),
					Used:      true, // Already used
				}
				mockStore.On("GetPasswordResetToken", "used-token").Return(resetToken, nil)
			},
			expectedError: ErrTokenUsed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			tt.setupMock(mockStore)

			cfg := &config.Config{
				PasswordResetTokenExpiry: 24,
				AppURL:                   "http://localhost:5173",
			}

			service := &AuthService{
				Store:        mockStore,
				Cfg:          cfg,
				EmailService: NewEmailService(cfg),
			}

			err := service.ResetPassword(tt.token, tt.newPassword)

			assert.Equal(t, tt.expectedError, err)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthService_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		updates       map[string]interface{}
		setupMock     func(*MockStore)
		expectedError error
	}{
		{
			name:   "successful user update",
			userID: "1",
			updates: map[string]interface{}{
				"username": "newusername",
				"email":    "newemail@example.com",
			},
			setupMock: func(mockStore *MockStore) {
				updatedUser := &store.User{
					ID:       "1",
					Username: "newusername",
					Email:    stringPtr("newemail@example.com"),
				}
				mockStore.On("UpdateUser", "1", map[string]interface{}{
					"username": "newusername",
					"email":    "newemail@example.com",
				}).Return(updatedUser, nil)
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: "999",
			updates: map[string]interface{}{
				"username": "newusername",
			},
			setupMock: func(mockStore *MockStore) {
				mockStore.On("UpdateUser", "999", map[string]interface{}{
					"username": "newusername",
				}).Return(nil, assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			tt.setupMock(mockStore)

			cfg := &config.Config{
				PasswordResetTokenExpiry: 24,
				AppURL:                   "http://localhost:5173",
			}

			service := &AuthService{
				Store:        mockStore,
				Cfg:          cfg,
				EmailService: NewEmailService(cfg),
			}

			user, err := service.UpdateUser(tt.userID, tt.updates)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthService_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		setupMock     func(*MockStore)
		expectedError error
	}{
		{
			name:   "successful user deletion",
			userID: "1",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("DeleteUser", "1").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: "999",
			setupMock: func(mockStore *MockStore) {
				mockStore.On("DeleteUser", "999").Return(assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			tt.setupMock(mockStore)

			cfg := &config.Config{
				PasswordResetTokenExpiry: 24,
				AppURL:                   "http://localhost:5173",
			}

			service := &AuthService{
				Store:        mockStore,
				Cfg:          cfg,
				EmailService: NewEmailService(cfg),
			}

			err := service.DeleteUser(tt.userID)

			assert.Equal(t, tt.expectedError, err)
			mockStore.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
