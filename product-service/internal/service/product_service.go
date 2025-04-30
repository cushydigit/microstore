package service

import (
	"github.com/cushydigit/microstore/porduct-service/internal/models"
	"github.com/cushydigit/microstore/porduct-service/internal/repository"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) Create(p *models.Product) error {
	return s.Repo.Create(p)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.Repo.GetAll()
}

func (s *ProductService) GetByID(id int64) (*models.Product, error) {
	return s.Repo.GetByID(id)
}

func (s *ProductService) Delete(id int64) error {
	return s.Repo.Delete(id)
}
