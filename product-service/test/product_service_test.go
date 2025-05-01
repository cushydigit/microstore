package test

import (
	"testing"

	"github.com/cushydigit/microstore/product-service/internal/models"
	"github.com/cushydigit/microstore/product-service/internal/repository"
	"github.com/cushydigit/microstore/product-service/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestProductService(t *testing.T) {
	repo := repository.NewInMemoryProductRepo()
	svc := service.NewProductService(repo)

	p := &models.Product{
		Name:        "Test",
		Description: "Test product",
		Price:       9.99,
	}

	// test empty projects
	ps, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, len(ps), 0)

	// test create
	err = svc.Create(p)
	assert.NoError(t, err)

	// test get by id
	got, err := svc.GetByID(int64(1))
	assert.NoError(t, err)
	assert.Equal(t, got.Name, p.Name)
	assert.Equal(t, got.ID, int64(1))

	// test delete
	err = svc.Delete(got.ID)
	assert.NoError(t, err)

	// test after delete
	got, err = svc.GetByID(got.ID)
	assert.NoError(t, err)
	assert.Nil(t, got)

}
