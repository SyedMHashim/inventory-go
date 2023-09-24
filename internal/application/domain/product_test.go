package domain_test

import (
	"testing"

	"example.com/apiserver/internal/application/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct_Success(t *testing.T) {
	product, err := domain.NewProduct("1", "chair", 200)
	assert.NoError(t, err)
	assert.Equal(t, "1", product.ID)
	assert.Equal(t, "chair", product.Name)
	assert.Equal(t, 200.00, product.Price)
}

func TestNewProduct_ReturnsError_InvalidProductName(t *testing.T) {
	_, err := domain.NewProduct("1", "", 200)
	assert.Error(t, err)
	assert.ErrorIs(t, domain.ErrInvalidProductName, err)
}

func TestNewProduct_ReturnsError_InvalidProductPrice(t *testing.T) {
	_, err := domain.NewProduct("1", "chair", 0)
	assert.Error(t, err)
	assert.ErrorIs(t, domain.ErrInvalidProductPrice, err)
}
