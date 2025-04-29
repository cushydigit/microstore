package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cushydigit/microstore/auth-service/internal/handler"
	"github.com/cushydigit/microstore/auth-service/internal/repository"
	"github.com/cushydigit/microstore/auth-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupRouter() http.Handler {
	repo := repository.NewInMemoryUserRepo()
	authService := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	return r
}

func TestRegisterAndLoginHandlerInMemory(t *testing.T) {
	router := setupRouter()

	// Register
	registerBody := map[string]string{
		"email":    "test@user.com",
		"password": "test123",
	}
	regJSON, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(regJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// login
	loginBody := map[string]string{
		"email":    "test@user.com",
		"password": "test123",
	}
	loginJSON, _ := json.Marshal(loginBody)

	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)

	var jsonResp map[string]any
	err := json.NewDecoder(resp.Body).Decode(&jsonResp)
	assert.NoError(t, err)

	assert.False(t, jsonResp["error"].(bool))
	assert.NotEmpty(t, jsonResp["data"].(map[string]any)["token"])

}
