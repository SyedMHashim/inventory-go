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

func TestGetProductService_Success_WhenGettingProduct(t *testing.T) {
	expectedProduct := domain.Product{
		ID:    "1",
		Name:  "Chair",
		Price: 200,
	}

	p := &out.MockGetProductPort{}
	p.MockGetProductByID(func(id string) (domain.Product, error) {
		return expectedProduct, nil
	})
	getProductService := service.NewGetProductService(p)
	actualProduct, err := getProductService.GetProduct(context.Background(), expectedProduct.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, actualProduct.ID)
	assert.Equal(t, expectedProduct.Name, actualProduct.Name)
	assert.Equal(t, expectedProduct.Price, actualProduct.Price)
}

func TestGetProductService_ReturnsError_WhenGettingProduct(t *testing.T) {
	p := &out.MockGetProductPort{}
	p.MockGetProductByID(func(id string) (domain.Product, error) {
		return domain.Product{}, errors.New("product does not exist")
	})
	getProductService := service.NewGetProductService(p)
	_, err := getProductService.GetProduct(context.Background(), "1")
	assert.Error(t, err)
}

func TestGetProductService_Success_WhenGettingProducts(t *testing.T) {
	expectedProducts := []domain.Product{
		{
			ID:    "1",
			Name:  "Chair",
			Price: 200,
		},
		{
			ID:    "2",
			Name:  "Table",
			Price: 500,
		},
	}

	p := &out.MockGetProductPort{}
	p.MockGetAllProducts(func() ([]domain.Product, error) {
		return expectedProducts, nil
	})
	getProductService := service.NewGetProductService(p)
	actualProducts, err := getProductService.GetProducts(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts[0].ID, actualProducts[0].ID)
	assert.Equal(t, expectedProducts[0].Name, actualProducts[0].Name)
	assert.Equal(t, expectedProducts[0].Price, actualProducts[0].Price)
	assert.Equal(t, expectedProducts[1].ID, actualProducts[1].ID)
	assert.Equal(t, expectedProducts[1].Name, actualProducts[1].Name)
	assert.Equal(t, expectedProducts[1].Price, actualProducts[1].Price)
}

func TestGetProductService_ReturnsError_WhenGettingProducts(t *testing.T) {
	p := &out.MockGetProductPort{}
	p.MockGetAllProducts(func() ([]domain.Product, error) {
		return []domain.Product{}, errors.New("error getting products")
	})
	getProductService := service.NewGetProductService(p)
	_, err := getProductService.GetProducts(context.Background())
	assert.Error(t, err)
}
