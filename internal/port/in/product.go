package in

import (
	"context"

	"example.com/apiserver/internal/application/domain"
)

type CreateProductUseCase interface {
	CreateProduct(ctx context.Context, name string, price float64) (string, error)
}

type GetProductUseCase interface {
	GetProduct(ctx context.Context, id string) (domain.Product, error)
	GetProducts(ctx context.Context) ([]domain.Product, error)
}

type UpdateProductUseCase interface {
	UpdateProduct(ctx context.Context, id, name string, price float64) error
}

type DeleteProductUseCase interface {
	DeleteProduct(ctx context.Context, id string) error
}
