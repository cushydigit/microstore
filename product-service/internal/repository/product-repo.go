package repository

import (
	"context"

	"github.com/cushydigit/microstore/shared/types"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*types.Product, error)
	Create(ctx context.Context, product *types.Product) error
	CreateBulk(ctx context.Context, product []*types.Product) error
	GetAll(ctx context.Context) ([]types.Product, error)
	Delete(ctx context.Context, id int64) error
	DeleteAll(ctx context.Context) error
}
