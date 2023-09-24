package out

import (
	"context"

	"example.com/apiserver/internal/application/domain"
)

type UpsertProductPort interface {
	InsertProduct(context.Context, domain.Product) error
	UpdateProduct(context.Context, domain.Product) error
}

type GetProductPort interface {
	GetProductByID(ctx context.Context, id string) (domain.Product, error)
	GetAllProducts(context.Context) ([]domain.Product, error)
}

type DeleteProductPort interface {
	DeleteProduct(ctx context.Context, id string) error
}
