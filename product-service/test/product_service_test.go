package test

import (
	"context"
	"testing"

	"github.com/cushydigit/microstore/product-service/internal/repository"
	"github.com/cushydigit/microstore/product-service/internal/service"
	"github.com/cushydigit/microstore/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestProductService(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewInMemoryProductRepo()
	svc := service.NewProductService(repo)

	p := &types.Product{
		Name:        "Test",
		Description: "Test product",
		Price:       9.99,
	}

	// test empty projects
	ps, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(ps), 0)

	// test create
	err = svc.Create(ctx, p)
	assert.NoError(t, err)

	// test get by id
	got, _, err := svc.GetByID(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, got.Name, p.Name)
	assert.Equal(t, got.ID, int64(1))

	// test delete
	err = svc.Delete(ctx, got.ID)
	assert.NoError(t, err)

	// test after delete
	got, _, err = svc.GetByID(ctx, got.ID)
	assert.NoError(t, err)
	assert.Nil(t, got)

}
