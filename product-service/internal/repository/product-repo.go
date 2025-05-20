package repository

import "github.com/cushydigit/microstore/shared/types"

type ProductRepository interface {
	GetByIDWithCache(id int64) (*types.Product, bool, error)
	GetByID(id int64) (*types.Product, error)
	Create(product *types.Product) error
	CreateBulk(product []types.Product) error
	GetAll() ([]types.Product, error)
	Delete(id int64) error
}
