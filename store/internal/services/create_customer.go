package services

import (
	"context"
	"log/slog"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type CreateCustomerReq struct {
	FirstName string
	LastName  string
	Email     *string
	Phone     *string
	UserID    string
}

type CreateCustomerResp = CreateUpdateResponse

type CreateCustomer struct {
	customerDAO dao.CustomerDAO
	eventBus    domain.EventBus
	nextID      domain.NextID
}

func NewCreateCustomer(customerDAO dao.CustomerDAO, eventBus domain.EventBus, nextID domain.NextID) *CreateCustomer {
	return &CreateCustomer{
		customerDAO: customerDAO,
		eventBus:    eventBus,
		nextID:      nextID,
	}
}

func (s *CreateCustomer) Exec(ctx context.Context, req CreateCustomerReq) (*CreateCustomerResp, error) {
	slog.InfoContext(ctx, "create customer started", "req", req)

	customer, err := domain.NewCustomer(s.nextID(), req.FirstName, req.LastName, req.Email, req.Phone, req.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "new customer failed", "error", err.Error())
		return nil, err
	}

	if err := s.customerDAO.Create(ctx, customer); err != nil {
		slog.ErrorContext(ctx, "create customer failed", "error", err.Error())
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, customer.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return nil, err
	}

	slog.InfoContext(ctx, "create customer finished", "customer_id", customer.GetID())

	return &CreateCustomerResp{ID: customer.GetID()}, nil
}
