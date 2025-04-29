package test

import (
	"testing"

	"github.com/cushydigit/microstore/auth-service/internal/repository"
	"github.com/cushydigit/microstore/auth-service/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLoginInMemory(t *testing.T) {
	repo := repository.NewInMemoryUserRepo()
	authService := service.NewAuthService(repo)

	email := "test@example.com"
	password := "secure123"

	// test registration
	user, err := authService.Register(email, password)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NotZero(t, user.ID)

	// test duplication
	_, err = authService.Register(email, password)
	assert.Error(t, err)

	// test login
	token, loggedUser, err := authService.Login(email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.ID, loggedUser.ID)
	assert.Equal(t, user.Email, loggedUser.Email)

	// test login failure
	_, _, err = authService.Login(email, "chert")
	assert.Error(t, err)

	// test login failure - unknown user
	_, _, err = authService.Login("chert@gmail.com", password)
	assert.Error(t, err)
}
