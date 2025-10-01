package services

import (
	"context"
	"log/slog"

	"ichibuy/order/internal/domain"
	"ichibuy/order/internal/domain/dao"
)

type CreateOrderReq struct {
	OrderLines []OrderLineReq
	UserID     string
}

type OrderLineReq struct {
	ProductID         string
	ProductName       string
	ProductStoreID    string
	Quantity          int
	UnitPriceAmount   int
	UnitPriceCurrency string
}

type CreateOrderResp struct {
	ID string `json:"id"`
}

type CreateOrder struct {
	orderDAO     dao.OrderDAO
	eventBus     domain.EventBus
	nextID       domain.NextID
	orderFactory *domain.OrderFactory
}

func NewCreateOrder(orderDAO dao.OrderDAO, eventBus domain.EventBus, nextID domain.NextID, orderFactory *domain.OrderFactory) *CreateOrder {
	return &CreateOrder{
		orderDAO:     orderDAO,
		eventBus:     eventBus,
		nextID:       nextID,
		orderFactory: orderFactory,
	}
}

func (s *CreateOrder) Exec(ctx context.Context, req CreateOrderReq) (*CreateOrderResp, error) {
	slog.InfoContext(ctx, "create order started", "req", req)

	orderLines, err := s.mapOrderLines(ctx, req)
	if err != nil {
		return nil, err
	}

	order, err := s.orderFactory.NewOrder(ctx, orderLines, req.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "new order failed", "error", err.Error())
		return nil, err
	}

	if err := s.orderDAO.Create(ctx, order); err != nil {
		slog.ErrorContext(ctx, "create customer failed", "error", err.Error())
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, order.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return nil, err
	}

	slog.InfoContext(ctx, "create order finished", "order_id", order.GetID())

	return &CreateOrderResp{ID: order.GetID()}, nil
}

func (s *CreateOrder) mapOrderLines(ctx context.Context, req CreateOrderReq) ([]domain.OrderLine, error) {
	orderLines := []domain.OrderLine{}
	for _, orderLineReq := range req.OrderLines {
		unitPrice, err := domain.NewMoney(orderLineReq.UnitPriceAmount, orderLineReq.UnitPriceCurrency)
		if err != nil {
			slog.ErrorContext(ctx, "new money failed", "error", err.Error())
			return nil, err
		}

		orderLine, err := domain.NewOrderLine(
			s.nextID(),
			orderLineReq.ProductID,
			orderLineReq.ProductName,
			orderLineReq.ProductStoreID,
			orderLineReq.Quantity,
			unitPrice,
		)
		if err != nil {
			slog.ErrorContext(ctx, "new order line failed", "error", err.Error())
			return nil, err
		}

		orderLines = append(orderLines, *orderLine)
	}
	return orderLines, nil
}
