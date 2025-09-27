package services

import (
	"context"
	"log/slog"

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
	slog.InfoContext(ctx, "delete product started", "req", req)
	product, err := s.productDAO.FindByPk(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "find product failed", "error", err.Error())
		return err
	}

	product.PrepareDelete()

	if err := s.productDAO.DeleteByPk(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "delete product failed", "error", err.Error())
		return err
	}

	if err := s.eventBus.Publish(ctx, product.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return err
	}

	slog.InfoContext(ctx, "delete product finished", "product_id", req.ID)
	return nil
}
