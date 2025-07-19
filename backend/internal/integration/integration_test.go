//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/user/votex-template/backend/internal/api"
	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/middleware"
	"github.com/user/votex-template/backend/internal/service"
	"github.com/user/votex-template/backend/internal/store"
)

var (
	testDB *sqlx.DB
	cfg    *config.Config
)

func TestMain(m *testing.M) {
	// Setup test environment
	setupTestEnv()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestEnv()

	os.Exit(code)
}

func setupTestEnv() {
	// Load test configuration
	cfg = &config.Config{
		Environment: "test",
		Port:        "8080",
		DBURL:       "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable",
		JWTSecret:   "test-secret-key",
		LogLevel:    "debug",
		CORSOrigins: []string{"http://localhost:5173"},
	}

	// Connect to test database
	var err error
	testDB, err = sqlx.Connect("postgres", cfg.DBURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	// Run migrations
	if err := runTestMigrations(); err != nil {
		panic(fmt.Sprintf("failed to run test migrations: %v", err))
	}
}

func cleanupTestEnv() {
	if testDB != nil {
		testDB.Close()
	}
}

func runTestMigrations() error {
	// Create tables for testing
	queries := []string{
		`CREATE TABLE IF NOT EXISTS "user" (
			id TEXT PRIMARY KEY,
			age INTEGER,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS session (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES "user"(id),
			expires_at TIMESTAMPTZ NOT NULL
		)`,
	}

	for _, query := range queries {
		if _, err := testDB.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func cleanupTestData() {
	testDB.Exec("DELETE FROM session")
	testDB.Exec("DELETE FROM \"user\"")
}

func setupTestServer() (*httptest.Server, *api.AuthHandler) {
	// Initialize store and service
	storeInstance := store.New(testDB)
	authService := service.NewAuthService(storeInstance, cfg)
	authHandler := api.NewAuthHandler(authService)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/auth/register":
			authHandler.Register(w, r)
		case "/api/auth/login":
			authHandler.Login(w, r)
		case "/api/auth/profile":
			// Add middleware for profile endpoint
			authMiddleware := middleware.NewAuthMiddleware(cfg)
			authMiddleware.Authenticate(http.HandlerFunc(authHandler.Profile)).ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	return server, authHandler
}

func TestIntegration_AuthFlow(t *testing.T) {
	cleanupTestData()
	defer cleanupTestData()

	server, _ := setupTestServer()
	defer server.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("Complete Auth Flow", func(t *testing.T) {
		// Step 1: Register a new user
		registerData := map[string]string{
			"username": "integrationtest",
			"password": "password123",
		}
		registerBody, _ := json.Marshal(registerData)

		registerReq, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewBuffer(registerBody))
		registerReq.Header.Set("Content-Type", "application/json")

		registerResp, err := client.Do(registerReq)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}
		defer registerResp.Body.Close()

		if registerResp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", registerResp.StatusCode)
		}

		var registerResult api.Response
		if err := json.NewDecoder(registerResp.Body).Decode(&registerResult); err != nil {
			t.Fatalf("failed to decode register response: %v", err)
		}

		if !registerResult.Success {
			t.Errorf("expected success response, got error: %s", registerResult.Error)
		}

		// Extract token from response
		authData, ok := registerResult.Data.(map[string]interface{})
		if !ok {
			t.Fatal("invalid response data format")
		}
		token, ok := authData["token"].(string)
		if !ok {
			t.Fatal("token not found in response")
		}

		// Step 2: Login with the same credentials
		loginData := map[string]string{
			"username": "integrationtest",
			"password": "password123",
		}
		loginBody, _ := json.Marshal(loginData)

		loginReq, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewBuffer(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")

		loginResp, err := client.Do(loginReq)
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}
		defer loginResp.Body.Close()

		if loginResp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", loginResp.StatusCode)
		}

		// Step 3: Access profile with token
		profileReq, _ := http.NewRequest("GET", server.URL+"/api/auth/profile", nil)
		profileReq.Header.Set("Authorization", "Bearer "+token)

		profileResp, err := client.Do(profileReq)
		if err != nil {
			t.Fatalf("failed to get profile: %v", err)
		}
		defer profileResp.Body.Close()

		if profileResp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", profileResp.StatusCode)
		}

		var profileResult api.Response
		if err := json.NewDecoder(profileResp.Body).Decode(&profileResult); err != nil {
			t.Fatalf("failed to decode profile response: %v", err)
		}

		if !profileResult.Success {
			t.Errorf("expected success response, got error: %s", profileResult.Error)
		}
	})

	t.Run("Invalid Login", func(t *testing.T) {
		loginData := map[string]string{
			"username": "integrationtest",
			"password": "wrongpassword",
		}
		loginBody, _ := json.Marshal(loginData)

		loginReq, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewBuffer(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")

		loginResp, err := client.Do(loginReq)
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}
		defer loginResp.Body.Close()

		if loginResp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", loginResp.StatusCode)
		}
	})

	t.Run("Unauthorized Profile Access", func(t *testing.T) {
		profileReq, _ := http.NewRequest("GET", server.URL+"/api/auth/profile", nil)
		// No Authorization header

		profileResp, err := client.Do(profileReq)
		if err != nil {
			t.Fatalf("failed to get profile: %v", err)
		}
		defer profileResp.Body.Close()

		if profileResp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", profileResp.StatusCode)
		}
	})
}

func TestIntegration_DatabaseOperations(t *testing.T) {
	cleanupTestData()
	defer cleanupTestData()

	storeInstance := store.New(testDB)

	t.Run("User Creation and Retrieval", func(t *testing.T) {
		// Create user
		userID := "test-user-123"
		username := "testuser"
		passwordHash := "hashedpassword"

		err := storeInstance.CreateUser(userID, username, passwordHash)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		// Retrieve user by ID
		user, err := storeInstance.GetUserByID(userID)
		if err != nil {
			t.Fatalf("failed to get user by ID: %v", err)
		}

		if user.ID != userID {
			t.Errorf("expected user ID %s, got %s", userID, user.ID)
		}

		if user.Username != username {
			t.Errorf("expected username %s, got %s", username, user.Username)
		}

		// Retrieve user by username
		userByUsername, err := storeInstance.GetUserByUsername(username)
		if err != nil {
			t.Fatalf("failed to get user by username: %v", err)
		}

		if userByUsername.ID != userID {
			t.Errorf("expected user ID %s, got %s", userID, userByUsername.ID)
		}
	})

	t.Run("Session Operations", func(t *testing.T) {
		sessionID := "test-session-123"
		userID := "test-user-456"
		expiresAt := "2024-12-31T23:59:59Z"

		// Create session
		err := storeInstance.CreateSession(sessionID, userID, expiresAt)
		if err != nil {
			t.Fatalf("failed to create session: %v", err)
		}

		// Get session
		session, err := storeInstance.GetSession(sessionID)
		if err != nil {
			t.Fatalf("failed to get session: %v", err)
		}

		if session.ID != sessionID {
			t.Errorf("expected session ID %s, got %s", sessionID, session.ID)
		}

		if session.UserID != userID {
			t.Errorf("expected user ID %s, got %s", userID, session.UserID)
		}

		// Delete session
		err = storeInstance.DeleteSession(sessionID)
		if err != nil {
			t.Fatalf("failed to delete session: %v", err)
		}

		// Verify session is deleted
		_, err = storeInstance.GetSession(sessionID)
		if err == nil {
			t.Error("expected error when getting deleted session")
		}
	})
}
