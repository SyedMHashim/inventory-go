package service

import (
	"context"
	"log"

	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/port/out"
)

type UpdateProductService struct {
	productPort out.UpsertProductPort
}

func NewUpdateProductService(
	productPort out.UpsertProductPort,
) *UpdateProductService {
	return &UpdateProductService{
		productPort: productPort,
	}
}

func (s *UpdateProductService) UpdateProduct(
	ctx context.Context,
	id string,
	name string,
	price float64,
) error {
	log.Println("Updating Product")
	product, err := domain.NewProduct(id, name, price)
	if err != nil {
		return err
	}
	if err := s.productPort.UpdateProduct(ctx, product); err != nil {
		return err
	}
	return nil
}
