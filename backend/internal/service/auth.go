package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// Error definitions
var (
	ErrUserExists         = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenNotFound      = errors.New("password reset token not found")
	ErrTokenExpired       = errors.New("password reset token expired")
	ErrTokenUsed          = errors.New("password reset token already used")
)

type User struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	Age       *int       `json:"age,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// AuthServiceInterface defines the interface for authentication operations
type AuthServiceInterface interface {
	Register(username, email, password string) (string, *User, error)
	Login(username, password string) (string, *User, error)
	GetUserByID(userID string) (*User, error)
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error
	UpdateUser(userID string, updates map[string]interface{}) (*User, error)
	DeleteUser(userID string) error
}

type AuthService struct {
	Store        store.StoreInterface
	Cfg          *config.Config
	EmailService *EmailService
}

func NewAuthService(s store.StoreInterface, cfg *config.Config) AuthServiceInterface {
	return &AuthService{
		Store:        s,
		Cfg:          cfg,
		EmailService: NewEmailService(cfg),
	}
}

func (s *AuthService) Register(username, email, password string) (string, *User, error) {
	// Check if user already exists by username
	existingUser, err := s.Store.GetUserByUsername(username)
	if err == nil && existingUser != nil {
		return "", nil, ErrUserExists
	}

	// Check if user already exists by email
	if email != "" {
		existingUser, err = s.Store.GetUserByEmail(email)
		if err == nil && existingUser != nil {
			return "", nil, ErrEmailExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	// Create user in database
	user := &User{
		ID:       generateID(),
		Username: username,
		Email:    &email,
	}

	err = s.Store.CreateUser(user.ID, username, email, string(hashedPassword))
	if err != nil {
		return "", nil, err
	}

	// Send welcome email
	if email != "" {
		go func() {
			if err := s.EmailService.SendWelcomeEmail(email, username); err != nil {
				// Log error but don't fail registration
				fmt.Printf("Failed to send welcome email: %v\n", err)
			}
		}()
	}

	token, err := s.generateToken(user.ID, username)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Login(username, password string) (string, *User, error) {
	// Get user from database
	dbUser, err := s.Store.GetUserByUsername(username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	user := &User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Age:       dbUser.Age,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	token, err := s.generateToken(user.ID, username)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GetUserByID(userID string) (*User, error) {
	dbUser, err := s.Store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Age:       dbUser.Age,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}

func (s *AuthService) RequestPasswordReset(email string) error {
	// Get user by email
	user, err := s.Store.GetUserByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate reset token
	token := generateSecureToken()
	expiresAt := time.Now().Add(time.Duration(s.Cfg.PasswordResetTokenExpiry) * time.Hour)

	// Store reset token
	err = s.Store.CreatePasswordResetToken(generateID(), user.ID, token, expiresAt)
	if err != nil {
		return err
	}

	// Send password reset email
	return s.EmailService.SendPasswordResetEmail(email, token)
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Get reset token
	resetToken, err := s.Store.GetPasswordResetToken(token)
	if err != nil {
		return ErrTokenNotFound
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return ErrTokenExpired
	}

	// Check if token is already used
	if resetToken.Used {
		return ErrTokenUsed
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	updates := map[string]interface{}{
		"password_hash": string(hashedPassword),
	}
	err = s.Store.UpdateUser(resetToken.UserID, updates)
	if err != nil {
		return err
	}

	// Mark token as used
	return s.Store.MarkPasswordResetTokenUsed(resetToken.ID)
}

func (s *AuthService) UpdateUser(userID string, updates map[string]interface{}) (*User, error) {
	// Check if user exists
	_, err := s.Store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update user
	err = s.Store.UpdateUser(userID, updates)
	if err != nil {
		return nil, err
	}

	// Get updated user
	return s.GetUserByID(userID)
}

func (s *AuthService) DeleteUser(userID string) error {
	// Check if user exists
	_, err := s.Store.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.Store.DeleteUser(userID)
}

func (s *AuthService) generateToken(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Cfg.JWTSecret))
}

// generateID creates a simple ID for demo purposes
// In production, use a proper UUID library
func generateID() string {
	return time.Now().Format("20060102150405")
}

// generateSecureToken creates a secure random token
func generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
