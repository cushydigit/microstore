package repository

import (
	"github.com/cushydigit/microstore/shared/types"
)

type UserRepository interface {
	FindByEmail(email string) (*types.User, error)
	Create(user *types.User) error
}
