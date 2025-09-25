package services

import (
	"context"
	"encoding/json"

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

	if err := s.productDAO.DeleteByPk(ctx, req.ID); err != nil {
		return err
	}

	eventData := domain.ProductEventData{
		ID:          product.GetID(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Active:      product.GetActive(),
		StoreID:     product.GetStoreID(),
		Images:      product.GetImages(),
		Prices:      product.GetPrices(),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	data, _ := json.Marshal(eventData)
	event := domain.Event{
		ID:        s.nextID(),
		Type:      domain.ProductDeleted,
		Data:      data,
		Timestamp: product.UpdatedAt,
	}

	s.eventBus.Publish(event)

	return nil
}
