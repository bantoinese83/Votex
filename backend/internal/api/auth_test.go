package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/votex-template/backend/internal/service"
	"github.com/user/votex-template/backend/internal/store"
)

// MockAuthService is a mock implementation for testing
type MockAuthService struct {
	registerFunc func(username, password string) (string, *service.User, error)
	loginFunc    func(username, password string) (string, *service.User, error)
	getUserFunc  func(userID string) (*service.User, error)
}

func (m *MockAuthService) Register(username, password string) (string, *service.User, error) {
	if m.registerFunc != nil {
		return m.registerFunc(username, password)
	}
	return "", nil, nil
}

func (m *MockAuthService) Login(username, password string) (string, *service.User, error) {
	if m.loginFunc != nil {
		return m.loginFunc(username, password)
	}
	return "", nil, nil
}

func (m *MockAuthService) GetUserByID(userID string) (*service.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(userID)
	}
	return nil, nil
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    AuthRequest
		mockRegister   func(username, password string) (string, *service.User, error)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid registration",
			requestBody: AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockRegister: func(username, password string) (string, *service.User, error) {
				return "jwt-token", &service.User{ID: "123", Username: username}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "username too short",
			requestBody: AuthRequest{
				Username: "ab",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "password too short",
			requestBody: AuthRequest{
				Username: "testuser",
				Password: "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "service error",
			requestBody: AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockRegister: func(username, password string) (string, *service.User, error) {
				return "", nil, service.ErrUserExists
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				registerFunc: tt.mockRegister,
			}

			handler := NewAuthHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Register(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response Response
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				if !response.Success {
					t.Errorf("expected success response, got error: %s", response.Error)
				}
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    AuthRequest
		mockLogin      func(username, password string) (string, *service.User, error)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid login",
			requestBody: AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockLogin: func(username, password string) (string, *service.User, error) {
				return "jwt-token", &service.User{ID: "123", Username: username}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "invalid credentials",
			requestBody: AuthRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockLogin: func(username, password string) (string, *service.User, error) {
				return "", nil, service.ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				loginFunc: tt.mockLogin,
			}

			handler := NewAuthHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthHandler_Profile(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockGetUser    func(userID string) (*service.User, error)
		expectedStatus int
	}{
		{
			name:   "valid profile request",
			userID: "123",
			mockGetUser: func(userID string) (*service.User, error) {
				return &service.User{ID: userID, Username: "testuser"}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "456",
			mockGetUser: func(userID string) (*service.User, error) {
				return nil, store.ErrUserNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				getUserFunc: tt.mockGetUser,
			}

			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest("GET", "/api/auth/profile", nil)
			// Add user ID to context (simulating middleware)
			ctx := req.Context()
			ctx = context.WithValue(ctx, "user_id", tt.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.Profile(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandleHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	HandleHealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("expected success response, got error: %s", response.Error)
	}
}
