package domain

import (
	"github.com/google/uuid"
)

func GenerateOrderID() string {
	return uuid.NewString()
}

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(id string, name string, price float64) (Product, error) {
	if name == "" {
		return Product{}, ErrInvalidProductName
	}
	if price == 0 {
		return Product{}, ErrInvalidProductPrice
	}
	return Product{
		ID:    id,
		Name:  name,
		Price: price,
	}, nil
}
