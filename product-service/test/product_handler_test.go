package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cushydigit/microstore/product-service/internal/handler"
	"github.com/cushydigit/microstore/product-service/internal/repository"
	"github.com/cushydigit/microstore/product-service/internal/service"
	"github.com/cushydigit/microstore/shared/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupRouter() http.Handler {
	repo := repository.NewInMemoryProductRepo()
	svc := service.NewProductService(repo, nil)
	h := handler.NewProductHandler(svc)

	r := chi.NewRouter()
	r.Post("/product", h.Create)
	r.Post("/product/bulk", h.Create)
	r.Get("/product", h.GetAll)
	r.Get("/product/{id}", h.GetByID)
	r.Delete("/product/{id}", h.Delete)

	return r
}

func TestProductHandler(t *testing.T) {
	r := setupRouter()

	// test empty products
	req := httptest.NewRequest(http.MethodGet, "/product", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// test create
	product := types.Product{
		Name:        "Apple",
		Price:       1.23,
		Description: "Red lebonanian apple",
	}

	body, _ := json.Marshal(product)
	req = httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// test get all after create
	req = httptest.NewRequest(http.MethodGet, "/product", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// test get by id
	req = httptest.NewRequest(http.MethodGet, "/product/1", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// test delete
	req = httptest.NewRequest(http.MethodDelete, "/product/1", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

}
