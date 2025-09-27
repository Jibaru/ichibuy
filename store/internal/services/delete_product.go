package services

import (
	"context"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type DeleteProductReq struct {
	ID string
}

type DeleteProduct struct {
	productDAO dao.ProductDAO
	eventBus   domain.EventBus
	nextID     domain.NextID
}

func NewDeleteProduct(productDAO dao.ProductDAO, eventBus domain.EventBus, nextID domain.NextID) *DeleteProduct {
	return &DeleteProduct{
		productDAO: productDAO,
		eventBus:   eventBus,
		nextID:     nextID,
	}
}

func (s *DeleteProduct) Exec(ctx context.Context, req DeleteProductReq) error {
	product, err := s.productDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return err
	}

	product.PrepareDelete()

	if err := s.productDAO.DeleteByPk(ctx, req.ID); err != nil {
		return err
	}

	if err := s.eventBus.Publish(ctx, product.PullEvents()...); err != nil {
		return err
	}

	return nil
}
