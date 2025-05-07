package handler

import (
	"net/http"

	"github.com/cushydigit/microstore/order-service/internal/service"
)

type OrderHandler struct {
	OrderService *service.OrderService
}

func NewOrderHandler(orderSevice *service.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderSevice}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *OrderHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {

}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {

}
