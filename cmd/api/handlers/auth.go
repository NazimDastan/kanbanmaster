package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input services.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.Name = middleware.SanitizeString(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" || input.Email == "" || input.Password == "" {
		writeError(w, "Name, email and password are required", http.StatusBadRequest)
		return
	}

	if len(input.Password) < 6 {
		writeError(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Register(input)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			writeError(w, "Email already in use", http.StatusConflict)
			return
		}
		writeError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, resp, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input services.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		writeError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(input)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) || errors.Is(err, services.ErrInvalidPassword) {
			writeError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		writeError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, resp, http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		writeError(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.RefreshToken(body.RefreshToken)
	if err != nil {
		writeError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	writeJSON(w, resp, http.StatusOK)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUser(userID.(string))
	if err != nil {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	writeJSON(w, user, http.StatusOK)
}

func (h *AuthHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var body struct {
		AvatarURL string `json:"avatarUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AvatarURL == "" {
		writeError(w, "avatarUrl is required", http.StatusBadRequest)
		return
	}

	// Limit size — base64 image shouldn't exceed ~2MB
	if len(body.AvatarURL) > 2*1024*1024 {
		writeError(w, "Avatar too large (max 2MB)", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateAvatar(userID, body.AvatarURL)
	if err != nil {
		writeError(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}
	writeJSON(w, user, http.StatusOK)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	body.Name = middleware.SanitizeString(body.Name)
	body.Email = strings.TrimSpace(body.Email)
	if body.Name == "" || body.Email == "" {
		writeError(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateProfile(userID, body.Name, body.Email)
	if err != nil {
		writeError(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}
	writeJSON(w, user, http.StatusOK)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CurrentPassword == "" || body.NewPassword == "" {
		writeError(w, "Current and new password are required", http.StatusBadRequest)
		return
	}

	if len(body.NewPassword) < 6 {
		writeError(w, "New password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	err := h.authService.ChangePassword(userID, body.CurrentPassword, body.NewPassword)
	if err != nil {
		if errors.Is(err, services.ErrInvalidPassword) {
			writeError(w, "Current password is incorrect", http.StatusUnauthorized)
			return
		}
		writeError(w, "Failed to change password", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"message": "Password changed"}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  status,
	})
}
