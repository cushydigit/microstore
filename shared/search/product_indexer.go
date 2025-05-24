package search

import (
	"context"

	"github.com/cushydigit/microstore/shared/types"
)

type ProductIndexer interface {
	IndexProduct(ctx context.Context, index string, p *types.Product) error
	IndexBulkProduct(ctx context.Context, index string, ps []*types.Product) error
	DeleteProduct(ctx context.Context, index string, id int64) error
	DeleteAllProducts(ctx context.Context, index string) error
	SearchProduct(ctx context.Context, index, query string) ([]*types.Product, error)
}
