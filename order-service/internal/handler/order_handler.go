package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cushydigit/microstore/order-service/internal/service"
	"github.com/cushydigit/microstore/shared/helpers"
	"github.com/cushydigit/microstore/shared/types"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	OrderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderService}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.CreateOrderRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.ErrorJSON(w, errors.New("Invalid request"))
		return
	}

	// get user ID from context
	userID := r.Context().Value("user_id").(int)
	order, err := h.OrderService.Create(userID, req.Items)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	payload := &types.OrderResponse{
		Error:   false,
		Message: "order created",
		Data:    *order,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)

}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("Invalid order ID"))
		return
	}

	order, err := h.OrderService.GetByID(id)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("Not found"), http.StatusNotFound)
		return
	}

	payload := types.OrderResponse{
		Error:   false,
		Message: "success",
		Data:    *order,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)

}

func (h *OrderHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {

	// get user ID from context
	userID := r.Context().Value("user_id").(int)
	orders, err := h.OrderService.GetAllByUserID(userID)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	payload := types.OrdersResponse{
		Error:   false,
		Message: "success",
		Data:    orders,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)

}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderService.GetAll()
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	payload := types.OrdersResponse{
		Error:   false,
		Message: "success",
		Data:    orders,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)

}
