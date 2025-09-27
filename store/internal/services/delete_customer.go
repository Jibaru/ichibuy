package services

import (
	"context"
	"log/slog"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type DeleteCustomerReq struct {
	ID string
}

type DeleteCustomer struct {
	customerDAO dao.CustomerDAO
	eventBus    domain.EventBus
	nextID      domain.NextID
}

func NewDeleteCustomer(customerDAO dao.CustomerDAO, eventBus domain.EventBus, nextID domain.NextID) *DeleteCustomer {
	return &DeleteCustomer{
		customerDAO: customerDAO,
		eventBus:    eventBus,
		nextID:      nextID,
	}
}

func (s *DeleteCustomer) Exec(ctx context.Context, req DeleteCustomerReq) error {
	slog.InfoContext(ctx, "delete customer started", "req", req)
	customer, err := s.customerDAO.FindByPk(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "find customer failed", "error", err.Error())
		return err
	}

	customer.PrepareDelete()

	if err := s.customerDAO.DeleteByPk(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "delete customer failed", "error", err.Error())
		return err
	}

	if err := s.eventBus.Publish(ctx, customer.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return err
	}

	slog.InfoContext(ctx, "delete customer finished", "customer_id", req.ID)
	return nil
}
