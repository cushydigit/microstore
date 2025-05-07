package repository

import "github.com/cushydigit/microstore/order-service/internal/models"

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id int64) (*models.Order, error)
	GetByUserID(userID int) ([]*models.Order, error)
	GetAll() ([]*models.Order, error)
}
