package service

import (
	"github.com/cushydigit/microstore/order-service/internal/models"
	"github.com/cushydigit/microstore/order-service/internal/repository"
)

type OrderService struct {
	Repo repository.OrderRepository
	ProductAPIURL string
}

func NewOrderSevice (repo repository.OrderRepository, productAPIURL string) *OrderService {
	return &OrderService{
		Repo: repo,
		ProductAPIURL: productAPIURL,
	}
}

func (s *OrderService) CreateOrder(userID int, items[]models.OrderItem) (*models.Order, error) {
	return nil, nil
}

func (s *OrderService) GetOrder(id int64) (*models.Order, error) {
	return s.Repo.GetByID(id)
}

func (s *OrderService) GetOrdersByUserID(id int) ([]*models.Order, error ) {
	return s.Repo.GetByUserID(id)
}

func (s *OrderService) GetAllOrders()([]*models.Order, error) {
	return s.Repo.GetAll() 
}
