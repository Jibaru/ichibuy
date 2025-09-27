package services

import (
	"context"
	"log/slog"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type UpdateCustomerReq struct {
	ID        string
	FirstName string
	LastName  string
	Email     *string
	Phone     *string
	UserID    string
}

type UpdateCustomer struct {
	customerDAO dao.CustomerDAO
	eventBus    domain.EventBus
	nextID      domain.NextID
}

func NewUpdateCustomer(customerDAO dao.CustomerDAO, eventBus domain.EventBus, nextID domain.NextID) *UpdateCustomer {
	return &UpdateCustomer{
		customerDAO: customerDAO,
		eventBus:    eventBus,
		nextID:      nextID,
	}
}

func (s *UpdateCustomer) Exec(ctx context.Context, req UpdateCustomerReq) error {
	slog.InfoContext(ctx, "update customer started", "req", req)
	customer, err := s.customerDAO.FindByPk(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "find customer failed", "error", err.Error())
		return err
	}

	if err := customer.Update(req.FirstName, req.LastName, req.Email, req.Phone, req.UserID); err != nil {
		slog.ErrorContext(ctx, "update customer domain failed", "error", err.Error())
		return err
	}

	if err := s.customerDAO.Update(ctx, customer); err != nil {
		slog.ErrorContext(ctx, "update customer failed", "error", err.Error())
		return err
	}

	if err := s.eventBus.Publish(ctx, customer.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return err
	}

	slog.InfoContext(ctx, "update customer finished", "customer_id", customer.GetID())
	return nil
}
