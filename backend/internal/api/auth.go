package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/user/votex-template/backend/internal/service"
)

type AuthHandler struct {
	Service   *service.AuthService
	Validator *validator.Validate
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Service:   s,
		Validator: validator.New(),
	}
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
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

	token, user, err := h.Service.Register(req.Username, req.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}{
			ID:       user.ID,
			Username: user.Username,
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
		WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}{
			ID:       user.ID,
			Username: user.Username,
		},
	}

	WriteSuccess(w, response)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r)
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
		ID       string `json:"id"`
		Username string `json:"username"`
	}{
		ID:       user.ID,
		Username: user.Username,
	}

	WriteSuccess(w, response)
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "healthy",
		Message: "Backend is running",
	}
	WriteSuccess(w, response)
}
