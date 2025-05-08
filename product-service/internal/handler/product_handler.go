package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cushydigit/microstore/shared/helpers"
	"github.com/cushydigit/microstore/shared/types"

	"github.com/cushydigit/microstore/product-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: s}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p types.Product
	if err := helpers.ReadJSON(w, r, &p); err != nil {
		helpers.ErrorJSON(w, errors.New("Invalid request"))
		return
	}

	if err := h.ProductService.Create(&p); err != nil {
		log.Printf("error: %v", err)
		helpers.ErrorJSON(w, errors.New("failed to create product"), http.StatusInternalServerError)
		return
	}

	payload := types.ProductResponse{
		Error:   false,
		Message: "Product created",
		Data:    p,
	}
	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.GetAll()
	if err != nil {
		helpers.ErrorJSON(w, errors.New("failed to fetch products"), http.StatusInternalServerError)
		return
	}

	payload := types.ProductsResponse{
		Error:   false,
		Message: "success",
		Data:    products,
	}
	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid product ID"))
		return
	}

	p, err := h.ProductService.GetByID(id)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("error fetching product"), http.StatusInternalServerError)
		return
	}

	if p == nil {
		helpers.ErrorJSON(w, errors.New("product not found"), http.StatusNotFound)
	}

	payload := types.ProductResponse{
		Error:   false,
		Message: "success",
		Data:    *p,
	}
	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid product ID"))
		return
	}

	if err := h.ProductService.Delete(id); err != nil {
		helpers.ErrorJSON(w, errors.New("product not found"), http.StatusNotFound)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: fmt.Sprintf("product with id %d deleted", id),
		Data:    nil,
	}
	helpers.WriteJSON(w, http.StatusOK, payload)
}
