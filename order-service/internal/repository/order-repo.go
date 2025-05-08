package repository

import "github.com/cushydigit/microstore/shared/types"

type OrderRepository interface {
	Create(order *types.Order) error
	GetByID(id int64) (*types.Order, error)
	GetByUserID(userID int) ([]types.Order, error)
	GetAll() ([]types.Order, error)
}
