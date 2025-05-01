package repository

import "github.com/cushydigit/microstore/product-service/internal/models"

type ProductRepository interface {
	GetByID(id int64) (*models.Product, error)
	Create(product *models.Product) error
	GetAll() ([]models.Product, error)
	Delete(id int64) error
}
