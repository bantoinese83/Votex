package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AuthService struct {
	Store store.StoreInterface
	Cfg   *config.Config
}

func NewAuthService(s store.StoreInterface, cfg *config.Config) *AuthService {
	return &AuthService{Store: s, Cfg: cfg}
}

func (s *AuthService) Register(username, password string) (string, *User, error) {
	// Check if user already exists
	existingUser, err := s.Store.GetUserByUsername(username)
	if err == nil && existingUser != nil {
		return "", nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	// Create user in database
	user := &User{
		ID:       generateID(),
		Username: username,
	}

	err = s.Store.CreateUser(user.ID, username, string(hashedPassword))
	if err != nil {
		return "", nil, err
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
		return "", nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	user := &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
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
		return nil, err
	}

	return &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
	}, nil
}

func (s *AuthService) generateToken(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Cfg.JWTSecret))
}

// generateID creates a simple ID for demo purposes
// In production, use a proper UUID library
func generateID() string {
	return time.Now().Format("20060102150405")
}
