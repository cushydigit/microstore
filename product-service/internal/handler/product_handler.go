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
	p, ok := r.Context().Value(types.ProductKey).(types.Product)
	if !ok {
		helpers.ErrorJSON(w, errors.New("product not found in context"), http.StatusInternalServerError)
		return
	}

	if err := h.ProductService.Create(r.Context(), &p); err != nil {
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

func (h *ProductHandler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var ps []*types.Product
	if err := helpers.ReadJSON(w, r, &ps); err != nil {
		helpers.ErrorJSON(w, errors.New("Invalid request"))
		return
	}

	if err := h.ProductService.CreateBulk(r.Context(), ps); err != nil {
		helpers.ErrorJSON(w, errors.New("failed to create products"), http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "Bulk products created",
		Data:    nil,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.GetAll(r.Context())
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

	p, cached, err := h.ProductService.GetByID(r.Context(), id)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("error fetching product"), http.StatusInternalServerError)
		return
	}

	if p == nil {
		helpers.ErrorJSON(w, errors.New("product not found"), http.StatusNotFound)
		return
	}

	if cached {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
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

	if err := h.ProductService.Delete(r.Context(), id); err != nil {
		helpers.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: fmt.Sprintf("product with id %d deleted", id),
		Data:    nil,
	}
	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *ProductHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.ProductService.DeleteAll(r.Context()); err != nil {
		helpers.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	payload := types.Response{
		Error:   false,
		Message: "many products deleted from DB",
		Data:    nil,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *ProductHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		helpers.ErrorJSON(w, errors.New("missing query paramater ?q="))
		return
	}

	results, err := h.ProductService.Search(r.Context(), query)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed search product: %v", err), http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error: false,
	}

	if len(results) == 0 {
		payload.Message = "There is no result for search query"
		payload.Data = []types.Product{}
	} else {
		payload.Message = fmt.Sprintf("total products found with searh query: %d", len(results))
		payload.Data = results
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}
