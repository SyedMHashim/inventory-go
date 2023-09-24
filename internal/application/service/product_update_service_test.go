package service_test

import (
	"context"
	"testing"

	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/application/service"
	"example.com/apiserver/internal/port/out"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProductService_Success_WhenUpdatingProduct(t *testing.T) {
	p := &out.MockUpsertProductPort{}
	p.MockUpdateProduct(func(product domain.Product) error {
		return nil
	})
	updateProductService := service.NewUpdateProductService(p)
	err := updateProductService.UpdateProduct(context.Background(), "1", "chair", 200)
	assert.NoError(t, err)
}

func TestUpdateProductService_ReturnsError_WhenUpdatingProduct(t *testing.T) {
	p := &out.MockUpsertProductPort{}
	p.MockUpdateProduct(func(product domain.Product) error {
		return nil
	})
	updateProductService := service.NewUpdateProductService(p)
	err := updateProductService.UpdateProduct(context.Background(), "1", "chair", 200)
	assert.NoError(t, err)
}
