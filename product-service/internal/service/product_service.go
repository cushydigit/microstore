package service

import (
	"context"
	"log"
	"time"

	"github.com/cushydigit/microstore/product-service/internal/repository"
	myredis "github.com/cushydigit/microstore/shared/redis"
	"github.com/cushydigit/microstore/shared/types"

	"github.com/cushydigit/microstore/shared/search"
)

type ProductService struct {
	Repo          repository.ProductRepository
	searchIndexer search.ProductIndexer
	index         string
}

func NewProductService(repo repository.ProductRepository, searchIndexer search.ProductIndexer) *ProductService {
	return &ProductService{
		Repo:          repo,
		searchIndexer: searchIndexer,
		index:         "products",
	}
}

func (s *ProductService) Create(ctx context.Context, p *types.Product) error {
	if err := s.Repo.Create(ctx, p); err != nil {
		return err
	}
	// Index in zincsearch
	return s.searchIndexer.IndexProduct(ctx, s.index, p)
}

func (s *ProductService) CreateBulk(ctx context.Context, ps []*types.Product) error {
	if err := s.Repo.CreateBulk(ctx, ps); err != nil {
		return err
	}
	return s.searchIndexer.IndexBulkProduct(ctx, s.index, ps)
}

func (s *ProductService) GetAll(ctx context.Context) ([]types.Product, error) {
	return s.Repo.GetAll(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*types.Product, bool, error) {
	// Try cache
	product, found, err := myredis.GetProductFromCache(ctx, id)
	if err != nil {
		log.Printf("cache error: %v", err)
	}
	if found {
		return product, true, nil
	}

	// fallback to DB
	product, err = s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, false, err
	}

	// check if product is not found in DB
	if product == nil {
		return nil, false, nil
	}

	// set in cache for next time
	go func(p *types.Product) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := myredis.SetProductToCache(ctx, product); err != nil {
			log.Printf("failed to set cache: %v", err)
		}
	}(product)

	return product, false, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}
	// Invalidate cache
	if err := myredis.DeleteProductFromCache(ctx, id); err != nil {
		log.Printf("failed to invalidate product cache: %v", err)
	}
	// Delete indexing
	return s.searchIndexer.DeleteProduct(ctx, s.index, id)
}

func (s *ProductService) DeleteAll(ctx context.Context) error {
	if err := s.Repo.DeleteAll(ctx); err != nil {
		return err
	}
	return s.searchIndexer.DeleteAllProducts(ctx, s.index)
}

func (s *ProductService) Search(ctx context.Context, query string) ([]*types.Product, error) {
	return s.searchIndexer.SearchProduct(ctx, s.index, query)
}
