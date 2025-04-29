package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cushydigit/microstore/auth-service/internal/helpers"
	"github.com/cushydigit/microstore/auth-service/internal/service"
	"github.com/cushydigit/microstore/auth-service/internal/types"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ErrorJSON(w, errors.New("Invalid request"))
		return
	}

	user, err := h.AuthService.Register(req.Email, req.Password)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "Registration successful",
		Data: map[string]any{
			"id":    user.ID,
			"email": user.Email,
		},
	}
	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ErrorJSON(w, errors.New("invalid request"))
		return
	}

	token, user, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "Login successful",
		Data: map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"token": token,
		},
	}
	helpers.WriteJSON(w, http.StatusAccepted, payload)
}
