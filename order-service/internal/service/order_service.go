package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cushydigit/microstore/order-service/internal/repository"
	"github.com/cushydigit/microstore/shared/types"
)

var productEndpoint = os.Getenv("PRODUCT_API_URL")

type OrderService struct {
	Repo          repository.OrderRepository
	ProductAPIURL string
}

func NewOrderService(repo repository.OrderRepository, productAPIURL string) *OrderService {
	return &OrderService{
		Repo:          repo,
		ProductAPIURL: productAPIURL,
	}
}

func (s *OrderService) Create(userID int, items []types.OrderItem) (*types.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("no items provided")
	}

	totalPrice := 0.0
	for _, item := range items {
		resp, err := http.Get(fmt.Sprintf("%s/%d", productEndpoint, item.ProductID))
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch product %d", item.ProductID)
		}
		defer resp.Body.Close()

		productResp := types.ProductResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
			return nil, errors.New("failed to read response from service")
		}

		if item.Quantity > productResp.Data.Stock {
			return nil, fmt.Errorf("product %d out of stock", item.ProductID)
		}

		totalPrice += productResp.Data.Price * float64(item.Quantity)
	}

	order := &types.Order{
		UserID:     userID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	if err := s.Repo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetByID(id int64) (*types.Order, error) {
	return s.Repo.GetByID(id)
}

func (s *OrderService) GetAllByUserID(id int) ([]types.Order, error) {
	return s.Repo.GetByUserID(id)
}

func (s *OrderService) GetAll() ([]types.Order, error) {
	return s.Repo.GetAll()
}
