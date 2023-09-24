package service

import (
	"context"
	"log"

	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/port/out"
)

type GetProductService struct {
	productPort out.GetProductPort
}

func NewGetProductService(
	productPort out.GetProductPort,
) *GetProductService {
	return &GetProductService{
		productPort: productPort,
	}
}

func (s *GetProductService) GetProduct(
	ctx context.Context,
	id string,
) (domain.Product, error) {
	log.Println("Getting Product")
	product, err := s.productPort.GetProductByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *GetProductService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	log.Println("Getting All Products")
	products, err := s.productPort.GetAllProducts(ctx)
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}
