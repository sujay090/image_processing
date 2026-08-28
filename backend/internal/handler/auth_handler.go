package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"image_processing_backend/internal/db"
	"image_processing_backend/internal/service"
)

type AuthHandler struct {
	store     *db.Store
	jwtSecret string
}

func NewAuthHandler(store *db.Store, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

// Utility to send JSON responses
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Utility to send JSON errors
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "Missing required fields or password too short")
		return
	}

	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	userID, err := h.store.CreateUser(r.Context(), req.Username, req.Email, hashedPassword)
	if err != nil {
		sendError(w, http.StatusConflict, "Username or email already exists")
		return
	}

	otp, err := service.GenerateOTP()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to generate OTP")
		return
	}

	expiresAt := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	err = h.store.CreateOTP(r.Context(), userID, otp, expiresAt)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save OTP")
		return
	}

	service.SendMockEmail(req.Email, otp)

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully. Please check your email for the OTP.",
		"user_id": userID,
	})
}

type VerifyRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || len(req.OTP) != 6 {
		sendError(w, http.StatusBadRequest, "Invalid email or OTP format")
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		sendError(w, http.StatusUnauthorized, "User not found")
		return
	}

	if user.IsVerified {
		sendError(w, http.StatusBadRequest, "User is already verified")
		return
	}

	latestOTP, err := h.store.GetLatestOTP(r.Context(), user.ID)
	if err != nil || latestOTP == nil {
		sendError(w, http.StatusBadRequest, "No OTP found")
		return
	}

	if latestOTP.Code != req.OTP || time.Now().After(latestOTP.ExpiresAt) {
		sendError(w, http.StatusUnauthorized, service.ErrInvalidOTP.Error())
		return
	}

	err = h.store.VerifyUser(r.Context(), user.ID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to verify user")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Email verified successfully"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		sendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !user.IsVerified {
		sendError(w, http.StatusForbidden, "Please verify your email before logging in")
		return
	}

	if !service.CheckPasswordHash(req.Password, user.PasswordHash) {
		sendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := service.GenerateJWT(user.ID, h.jwtSecret)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
