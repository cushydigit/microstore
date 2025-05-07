package service

import (
	"errors"

	"github.com/cushydigit/microstore/auth-service/internal/repository"
	"github.com/cushydigit/microstore/auth-service/internal/utils"
	"github.com/cushydigit/microstore/shared/types"
)

type AuthService struct {
	Repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(email, password string) (*types.User, error) {
	if existing, _ := s.Repo.FindByEmail(email); existing != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, _ := utils.HashPassword(password)
	user := &types.User{
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *types.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
