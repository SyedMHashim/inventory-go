package in

import (
	"context"

	"example.com/apiserver/internal/application/domain"
)

type MockCreateProductUseCase struct {
	createProduct func(name string, price float64) (string, error)
}

func (uc *MockCreateProductUseCase) MockCreateProduct(createProduct func(name string, price float64) (string, error)) {
	uc.createProduct = createProduct
}

func (uc *MockCreateProductUseCase) CreateProduct(ctx context.Context, name string, price float64) (string, error) {
	return uc.createProduct(name, price)
}

type MockDeleteProductUseCase struct {
	deleteProduct func(id string) error
}

func (uc *MockDeleteProductUseCase) MockDeleteProduct(deleteProduct func(id string) error) {
	uc.deleteProduct = deleteProduct
}

func (uc *MockDeleteProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	return uc.deleteProduct(id)
}

type MockGetProductUseCase struct {
	getProduct  func(id string) (domain.Product, error)
	getProducts func() ([]domain.Product, error)
}

func (uc *MockGetProductUseCase) MockGetProduct(getProduct func(id string) (domain.Product, error)) {
	uc.getProduct = getProduct
}

func (uc *MockGetProductUseCase) MockGetProducts(getProducts func() ([]domain.Product, error)) {
	uc.getProducts = getProducts
}

func (uc *MockGetProductUseCase) GetProduct(ctx context.Context, id string) (domain.Product, error) {
	return uc.getProduct(id)
}

func (uc *MockGetProductUseCase) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return uc.getProducts()
}

type MockUpdateProductUseCase struct {
	updateProduct func(id, name string, price float64) error
}

func (uc *MockUpdateProductUseCase) MockUpdateProduct(updateProduct func(id, name string, price float64) error) {
	uc.updateProduct = updateProduct
}

func (uc *MockUpdateProductUseCase) UpdateProduct(ctx context.Context, id, name string, price float64) error {
	return uc.updateProduct(id, name, price)
}
