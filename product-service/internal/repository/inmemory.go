package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/cushydigit/microstore/shared/types"
)

type InMemoryProductRepo struct {
	products map[int64]*types.Product
	mu       sync.Mutex
	nextID   int64
}

func NewInMemoryProductRepo() ProductRepository {
	return &InMemoryProductRepo{
		products: make(map[int64]*types.Product),
		nextID:   1,
	}

}

func (r *InMemoryProductRepo) GetByID(ctx context.Context, id int64) (*types.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[id]
	if !ok {
		return nil, nil
	}

	return product, nil
}

func (r *InMemoryProductRepo) Create(ctx context.Context, p *types.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.products[p.ID] = p
	r.nextID++

	return nil
}

func (r *InMemoryProductRepo) CreateBulk(ctx context.Context, ps []*types.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range ps {
		p.ID = r.nextID
		r.products[p.ID] = p
		r.nextID++
	}
	return nil
}

func (r *InMemoryProductRepo) GetAll(ctx context.Context) ([]types.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var products []types.Product
	for _, p := range r.products {
		products = append(products, *p)
	}

	return products, nil
}

func (r *InMemoryProductRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)

	return nil
}

func (r *InMemoryProductRepo) DeleteAll(ctx context.Context) error {
	return nil
}
