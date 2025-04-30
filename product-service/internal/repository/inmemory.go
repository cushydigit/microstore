package repository

import (
	"errors"
	"sync"

	"github.com/cushydigit/microstore/porduct-service/internal/models"
)

type InMemoryProductRepo struct {
	products map[int64]*models.Product
	mu       sync.Mutex
	nextID   int64
}

func NewInMemoryProductRepo() ProductRepository {
	return &InMemoryProductRepo{
		products: make(map[int64]*models.Product),
		nextID:   1,
	}
}

func (r *InMemoryProductRepo) GetByID(id int64) (*models.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[id]
	if !ok {
		return nil, nil
	}

	return product, nil
}

func (r *InMemoryProductRepo) Create(p *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.products[p.ID] = p
	r.nextID++

	return nil
}

func (r *InMemoryProductRepo) GetAll() ([]models.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var products []models.Product
	for _, p := range r.products {
		products = append(products, *p)
	}

	return products, nil
}

func (r *InMemoryProductRepo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)

	return nil
}
