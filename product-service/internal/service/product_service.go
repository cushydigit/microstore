package service

import (
	"github.com/cushydigit/microstore/product-service/internal/repository"
	"github.com/cushydigit/microstore/shared/types"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) Create(p *types.Product) error {
	return s.Repo.Create(p)
}

func (s *ProductService) GetAll() ([]types.Product, error) {
	return s.Repo.GetAll()
}

func (s *ProductService) GetByID(id int64) (*types.Product, error) {
	return s.Repo.GetByID(id)
}

func (s *ProductService) Delete(id int64) error {
	return s.Repo.Delete(id)
}
