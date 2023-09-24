package service

import (
	"context"
	"log"

	"example.com/apiserver/internal/port/out"
)

type DeleteProductService struct {
	productPort out.DeleteProductPort
}

func NewDeleteProductService(productPort out.DeleteProductPort) *DeleteProductService {
	return &DeleteProductService{
		productPort: productPort,
	}
}

func (s *DeleteProductService) DeleteProduct(
	ctx context.Context,
	id string,
) error {
	log.Println("Deleting Product")
	if err := s.productPort.DeleteProduct(ctx, id); err != nil {
		return err
	}
	return nil
}
