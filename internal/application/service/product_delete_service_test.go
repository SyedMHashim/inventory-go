package service_test

import (
	"context"
	"errors"
	"testing"

	"example.com/apiserver/internal/application/service"
	"example.com/apiserver/internal/port/out"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProductService_Success_WhenDeletingProduct(t *testing.T) {
	p := &out.MockDeleteProductPort{}
	p.MockDeleteProduct(func(id string) error {
		return nil
	})
	deleteProductService := service.NewDeleteProductService(p)
	err := deleteProductService.DeleteProduct(context.Background(), "1")
	assert.NoError(t, err)
}

func TestDeleteProductService_ReturnsError_WhenDeletingProduct(t *testing.T) {
	p := &out.MockDeleteProductPort{}
	p.MockDeleteProduct(func(id string) error {
		return errors.New("error encountered")
	})
	deleteProductService := service.NewDeleteProductService(p)
	err := deleteProductService.DeleteProduct(context.Background(), "1")
	assert.Error(t, err)
}
