package service

import (
	"context"
	"log"

	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/port/out"
)

type CreateProductService struct {
	generateProductID func() string
	productPort       out.UpsertProductPort
}

func NewCreateProductService(
	generateProductID func() string,
	productPort out.UpsertProductPort,
) *CreateProductService {
	return &CreateProductService{
		generateProductID: generateProductID,
		productPort:       productPort,
	}
}

func (s *CreateProductService) CreateProduct(
	ctx context.Context,
	name string,
	price float64,
) (string, error) {
	log.Println("Creating Product")
	id := s.generateProductID()
	product, err := domain.NewProduct(id, name, price)
	if err != nil {
		return "", err
	}
	if err := s.productPort.InsertProduct(ctx, product); err != nil {
		return "", err
	}
	return id, nil
}
