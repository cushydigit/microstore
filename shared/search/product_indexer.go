package search

import (
	"context"

	"github.com/cushydigit/microstore/shared/types"
)

type ProductIndexer interface {
	IndexProduct(ctx context.Context, index string, p *types.Product) error
	DeleteProduct(ctx context.Context, index string, id int64) error
	SearchProduct(ctx context.Context, query string) ([]*types.Product, error)
}
