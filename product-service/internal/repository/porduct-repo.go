package repository

import "github.com/cushydigit/microstore/porduct-service/internal/models"

type ProductRepository interface {
	GetByID(id int64) (*models.Porduct, error)
	Create(product *models.Porduct) error
	GetAll() ([]*models.Porduct, error)
	Delete(id int64) error
}	
