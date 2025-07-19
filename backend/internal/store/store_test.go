package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	cleanup := func() {
		sqlxDB.Close()
	}

	return sqlxDB, mock, cleanup
}

func TestStore_CreateUser(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	store := New(db)

	tests := []struct {
		name         string
		id           string
		username     string
		passwordHash string
		expectError  bool
		setupMock    func()
	}{
		{
			name:         "successful user creation",
			id:           "123",
			username:     "testuser",
			passwordHash: "hashedpassword",
			expectError:  false,
			setupMock: func() {
				mock.ExpectExec("INSERT INTO \"user\"").
					WithArgs("123", "testuser", "hashedpassword").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:         "database error",
			id:           "123",
			username:     "testuser",
			passwordHash: "hashedpassword",
			expectError:  true,
			setupMock: func() {
				mock.ExpectExec("INSERT INTO \"user\"").
					WithArgs("123", "testuser", "hashedpassword").
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := store.CreateUser(tt.id, tt.username, tt.passwordHash)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestStore_GetUserByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	store := New(db)

	tests := []struct {
		name        string
		id          string
		expectError bool
		setupMock   func()
	}{
		{
			name:        "user found",
			id:          "123",
			expectError: false,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "age"}).
					AddRow("123", "testuser", "hashedpassword", nil)
				mock.ExpectQuery("SELECT id, username, password_hash, age FROM \"user\" WHERE id = \\$1").
					WithArgs("123").
					WillReturnRows(rows)
			},
		},
		{
			name:        "user not found",
			id:          "456",
			expectError: true,
			setupMock: func() {
				mock.ExpectQuery("SELECT id, username, password_hash, age FROM \"user\" WHERE id = \\$1").
					WithArgs("456").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			user, err := store.GetUserByID(tt.id)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != ErrUserNotFound {
					t.Errorf("expected ErrUserNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("expected user but got nil")
				}
				if user.ID != tt.id {
					t.Errorf("expected user ID %s, got %s", tt.id, user.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestStore_GetUserByUsername(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	store := New(db)

	tests := []struct {
		name        string
		username    string
		expectError bool
		setupMock   func()
	}{
		{
			name:        "user found",
			username:    "testuser",
			expectError: false,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "age"}).
					AddRow("123", "testuser", "hashedpassword", nil)
				mock.ExpectQuery("SELECT id, username, password_hash, age FROM \"user\" WHERE username = \\$1").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
		},
		{
			name:        "user not found",
			username:    "nonexistent",
			expectError: true,
			setupMock: func() {
				mock.ExpectQuery("SELECT id, username, password_hash, age FROM \"user\" WHERE username = \\$1").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			user, err := store.GetUserByUsername(tt.username)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != ErrUserNotFound {
					t.Errorf("expected ErrUserNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("expected user but got nil")
				}
				if user.Username != tt.username {
					t.Errorf("expected username %s, got %s", tt.username, user.Username)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestStore_CreateSession(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	store := New(db)

	tests := []struct {
		name        string
		id          string
		userID      string
		expiresAt   string
		expectError bool
		setupMock   func()
	}{
		{
			name:        "successful session creation",
			id:          "session123",
			userID:      "user123",
			expiresAt:   "2024-01-01T00:00:00Z",
			expectError: false,
			setupMock: func() {
				mock.ExpectExec("INSERT INTO session").
					WithArgs("session123", "user123", "2024-01-01T00:00:00Z").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:        "database error",
			id:          "session123",
			userID:      "user123",
			expiresAt:   "2024-01-01T00:00:00Z",
			expectError: true,
			setupMock: func() {
				mock.ExpectExec("INSERT INTO session").
					WithArgs("session123", "user123", "2024-01-01T00:00:00Z").
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := store.CreateSession(tt.id, tt.userID, tt.expiresAt)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestMockStore(t *testing.T) {
	mockStore := &MockStore{}

	t.Run("CreateUser", func(t *testing.T) {
		err := mockStore.CreateUser("123", "testuser", "hashedpassword")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("GetUserByID", func(t *testing.T) {
		user, err := mockStore.GetUserByID("123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user == nil {
			t.Errorf("expected user but got nil")
		}
		if user.ID != "123" {
			t.Errorf("expected user ID 123, got %s", user.ID)
		}
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		user, err := mockStore.GetUserByUsername("testuser")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user == nil {
			t.Errorf("expected user but got nil")
		}
		if user.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", user.Username)
		}
	})

	t.Run("CreateSession", func(t *testing.T) {
		err := mockStore.CreateSession("session123", "user123", "2024-01-01T00:00:00Z")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("GetSession", func(t *testing.T) {
		session, err := mockStore.GetSession("session123")
		if err == nil {
			t.Errorf("expected error but got none")
		}
		if session != nil {
			t.Errorf("expected nil session but got %v", session)
		}
	})

	t.Run("DeleteSession", func(t *testing.T) {
		err := mockStore.DeleteSession("session123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
