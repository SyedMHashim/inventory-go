package service_test

import (
	"context"
	"errors"
	"testing"

	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/application/service"
	"example.com/apiserver/internal/port/out"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductService_Success_WhenCreatingProduct(t *testing.T) {
	p := &out.MockUpsertProductPort{}
	p.MockInsertProduct(func(product domain.Product) error {
		return nil
	})
	createProductService := service.NewCreateProductService(func() string { return "1" }, p)
	id, err := createProductService.CreateProduct(context.Background(), "chair", 200)
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
}

func TestCreateProductService_ReturnsError_WhenCreatingProduct(t *testing.T) {
	p := &out.MockUpsertProductPort{}
	p.MockInsertProduct(func(product domain.Product) error {
		return errors.New("error encountered")
	})
	createProductService := service.NewCreateProductService(func() string { return "1" }, p)
	_, err := createProductService.CreateProduct(context.Background(), "chair", 200)
	assert.Error(t, err)
}
