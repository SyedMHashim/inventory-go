package out

import (
	"context"

	"example.com/apiserver/internal/application/domain"
)

type MockUpsertProductPort struct {
	insertProduct func(domain.Product) error
	updateProduct func(domain.Product) error
}

func (p *MockUpsertProductPort) MockInsertProduct(insertProduct func(product domain.Product) error) {
	p.insertProduct = insertProduct
}

func (p *MockUpsertProductPort) InsertProduct(ctx context.Context, product domain.Product) error {
	return p.insertProduct(product)
}

func (p *MockUpsertProductPort) MockUpdateProduct(updateProduct func(product domain.Product) error) {
	p.updateProduct = updateProduct
}

func (p *MockUpsertProductPort) UpdateProduct(ctx context.Context, product domain.Product) error {
	return p.updateProduct(product)
}

type MockGetProductPort struct {
	getProductByID func(id string) (domain.Product, error)
	getAllProducts func() ([]domain.Product, error)
}

func (p *MockGetProductPort) MockGetProductByID(getProductByID func(id string) (domain.Product, error)) {
	p.getProductByID = getProductByID
}

func (p *MockGetProductPort) GetProductByID(ctx context.Context, id string) (domain.Product, error) {
	return p.getProductByID(id)
}

func (p *MockGetProductPort) MockGetAllProducts(getAllProducts func() ([]domain.Product, error)) {
	p.getAllProducts = getAllProducts
}

func (p *MockGetProductPort) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return p.getAllProducts()
}

type MockDeleteProductPort struct {
	deleteProduct func(id string) error
}

func (p *MockDeleteProductPort) MockDeleteProduct(deleteProduct func(id string) error) {
	p.deleteProduct = deleteProduct
}

func (p *MockDeleteProductPort) DeleteProduct(ctx context.Context, id string) error {
	return p.deleteProduct(id)
}
