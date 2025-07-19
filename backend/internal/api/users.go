package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/user/votex-template/backend/internal/middleware"
	"github.com/user/votex-template/backend/internal/service"
)

type UserHandler struct {
	Service   service.AuthServiceInterface
	Validator *validator.Validate
}

func NewUserHandler(s service.AuthServiceInterface) *UserHandler {
	return &UserHandler{
		Service:   s,
		Validator: validator.New(),
	}
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type UserResponse struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Email     *string `json:"email,omitempty"`
	Age       *int    `json:"age,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
}

type UserUpdateAdminRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=32"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Age      *int    `json:"age" validate:"omitempty,min=0,max=150"`
}

// ListUsers handles GET /api/users - list all users with pagination
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// For now, we'll implement a simple list without search
	// In a real implementation, you'd add search functionality to the store
	_ = search

	// Get users from service (this would need to be implemented in the service layer)
	// For now, we'll return a placeholder response
	users := []UserResponse{
		{
			ID:       "1",
			Username: "admin",
			Email:    nil,
			Age:      nil,
		},
	}

	response := UserListResponse{
		Users: users,
		Total: len(users),
		Page:  page,
		Limit: limit,
	}

	WriteSuccess(w, response)
}

// GetUser handles GET /api/users/{id} - get a specific user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Check if user is requesting their own profile or is admin
	requestingUserID, ok := middleware.GetUserID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// For now, only allow users to view their own profile
	// In a real implementation, you'd check for admin role
	if requestingUserID != userID {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			WriteError(w, http.StatusNotFound, "User not found")
		default:
			WriteError(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())
		}
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	}

	WriteSuccess(w, response)
}

// UpdateUser handles PUT /api/users/{id} - update a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Check if user is updating their own profile or is admin
	requestingUserID, ok := middleware.GetUserID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// For now, only allow users to update their own profile
	// In a real implementation, you'd check for admin role
	if requestingUserID != userID {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req UserUpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteValidationError(w, err)
		return
	}

	// Convert request to updates map
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
		switch err {
		case service.ErrUserNotFound:
			WriteError(w, http.StatusNotFound, "User not found")
		case service.ErrUserExists:
			WriteError(w, http.StatusConflict, "Username already exists")
		case service.ErrEmailExists:
			WriteError(w, http.StatusConflict, "Email already exists")
		default:
			WriteError(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		}
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	}

	WriteSuccess(w, response)
}

// DeleteUser handles DELETE /api/users/{id} - delete a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Check if user is deleting their own account or is admin
	requestingUserID, ok := middleware.GetUserID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// For now, only allow users to delete their own account
	// In a real implementation, you'd check for admin role
	if requestingUserID != userID {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	err := h.Service.DeleteUser(userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			WriteError(w, http.StatusNotFound, "User not found")
		default:
			WriteError(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		}
		return
	}

	WriteSuccess(w, map[string]string{
		"message": "User deleted successfully",
	})
}
