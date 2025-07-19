package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/user/votex-template/backend/internal/middleware"
	"github.com/user/votex-template/backend/internal/service"
)

type AuthHandler struct {
	Service   service.AuthServiceInterface
	Validator *validator.Validate
}

func NewAuthHandler(s service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		Service:   s,
		Validator: validator.New(),
	}
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        string  `json:"id"`
		Username  string  `json:"username"`
		Email     *string `json:"email,omitempty"`
		Age       *int    `json:"age,omitempty"`
		CreatedAt *string `json:"created_at,omitempty"`
		UpdatedAt *string `json:"updated_at,omitempty"`
	} `json:"user"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserUpdateRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=32"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Age      *int    `json:"age" validate:"omitempty,min=0,max=150"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteValidationError(w, err)
		return
	}

	token, user, err := h.Service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			WriteError(w, http.StatusConflict, "Username already exists")
		case service.ErrEmailExists:
			WriteError(w, http.StatusConflict, "Email already exists")
		default:
			WriteError(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		}
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID        string  `json:"id"`
			Username  string  `json:"username"`
			Email     *string `json:"email,omitempty"`
			Age       *int    `json:"age,omitempty"`
			CreatedAt *string `json:"created_at,omitempty"`
			UpdatedAt *string `json:"updated_at,omitempty"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Age:      user.Age,
		},
	}

	WriteSuccess(w, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteValidationError(w, err)
		return
	}

	token, user, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		default:
			WriteError(w, http.StatusInternalServerError, "Login failed: "+err.Error())
		}
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID        string  `json:"id"`
			Username  string  `json:"username"`
			Email     *string `json:"email,omitempty"`
			Age       *int    `json:"age,omitempty"`
			CreatedAt *string `json:"created_at,omitempty"`
			UpdatedAt *string `json:"updated_at,omitempty"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Age:      user.Age,
		},
	}

	WriteSuccess(w, response)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	response := struct {
		ID        string  `json:"id"`
		Username  string  `json:"username"`
		Email     *string `json:"email,omitempty"`
		Age       *int    `json:"age,omitempty"`
		CreatedAt *string `json:"created_at,omitempty"`
		UpdatedAt *string `json:"updated_at,omitempty"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	}

	WriteSuccess(w, response)
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteValidationError(w, err)
		return
	}

	err := h.Service.RequestPasswordReset(req.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to send password reset email: "+err.Error())
		return
	}

	WriteSuccess(w, map[string]string{
		"message": "If the email exists, a password reset link has been sent",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		WriteError(w, http.StatusBadRequest, "Token is required")
		return
	}

	var req PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteValidationError(w, err)
		return
	}

	err := h.Service.ResetPassword(token, req.Password)
	if err != nil {
		switch err {
		case service.ErrTokenNotFound:
			WriteError(w, http.StatusNotFound, "Invalid or expired token")
		case service.ErrTokenExpired:
			WriteError(w, http.StatusBadRequest, "Token has expired")
		case service.ErrTokenUsed:
			WriteError(w, http.StatusBadRequest, "Token has already been used")
		default:
			WriteError(w, http.StatusInternalServerError, "Failed to reset password: "+err.Error())
		}
		return
	}

	WriteSuccess(w, map[string]string{
		"message": "Password reset successfully",
	})
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteValidationError(w, err)
		return
	}

	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}

	user, err := h.Service.UpdateUser(userID, updates)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update profile: "+err.Error())
		return
	}

	response := struct {
		ID        string  `json:"id"`
		Username  string  `json:"username"`
		Email     *string `json:"email,omitempty"`
		Age       *int    `json:"age,omitempty"`
		CreatedAt *string `json:"created_at,omitempty"`
		UpdatedAt *string `json:"updated_at,omitempty"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	}

	WriteSuccess(w, response)
}

func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.Service.DeleteUser(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete account: "+err.Error())
		return
	}

	WriteSuccess(w, map[string]string{
		"message": "Account deleted successfully",
	})
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
		Version   string `json:"version"`
	}{
		Status:    "healthy",
		Message:   "Backend is running",
		Timestamp: "2024-01-01T00:00:00Z",
		Version:   "1.0.0",
	}
	WriteSuccess(w, response)
}
