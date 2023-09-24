package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	web "example.com/apiserver/internal/adapter/in/http"
	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/port/in"
	"github.com/stretchr/testify/assert"
)

func TestGetProductHandler_Success_WhenGettingProduct(t *testing.T) {
	product := domain.Product{
		ID:    "1",
		Name:  "Chair",
		Price: 200,
	}

	uc := &in.MockGetProductUseCase{}
	uc.MockGetProduct(func(id string) (domain.Product, error) {
		return product, nil
	})
	getProductHandler := web.NewGetProductHandler(uc)

	req, err := http.NewRequest(http.MethodGet, "/product/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductHandler.HandleGetProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	b, _ := json.Marshal(product)
	assert.JSONEq(t, fmt.Sprintf(
		`{"message": "Product retrieved successfully", "data": {"product": %s}}`, string(b),
	), rr.Body.String())
}

func TestGetProductHandler_ReturnsError_WhenGettingProduct(t *testing.T) {
	uc := &in.MockGetProductUseCase{}
	uc.MockGetProduct(func(id string) (domain.Product, error) {
		return domain.Product{}, errors.New("product does not exist")
	})
	getProductHandler := web.NewGetProductHandler(uc)

	req, err := http.NewRequest(http.MethodGet, "/product/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductHandler.HandleGetProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{
			"error":   "Something wrong happened at our end",
			"details": "Please try again after a while. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues"
		}`, rr.Body.String())
}

func TestGetProductHandler_Success_WhenGettingProducts(t *testing.T) {
	products := []domain.Product{
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

	uc := &in.MockGetProductUseCase{}
	uc.MockGetProducts(func() ([]domain.Product, error) {
		return products, nil
	})
	getProductHandler := web.NewGetProductHandler(uc)

	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductHandler.HandleGetProducts)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	b, _ := json.Marshal(products)
	assert.JSONEq(t, fmt.Sprintf(
		`{"message": "Products retrieved successfully", "data": {"products": %s}}`, string(b),
	), rr.Body.String())
}

func TestGetProductHandler_ReturnsError_WhenGettingProducts(t *testing.T) {
	uc := &in.MockGetProductUseCase{}
	uc.MockGetProducts(func() ([]domain.Product, error) {
		return []domain.Product{}, errors.New("error getting products")
	})
	getProductHandler := web.NewGetProductHandler(uc)

	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductHandler.HandleGetProducts)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{
			"error":   "Something wrong happened at our end",
			"details": "Please try again after a while. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues"
		}`, rr.Body.String())
}
