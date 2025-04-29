package repository

import "github.com/cushydigit/microstore/auth-service/internal/models"

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}
