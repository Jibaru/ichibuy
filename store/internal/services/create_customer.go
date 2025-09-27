package services

import (
	"context"

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
	customer, err := domain.NewCustomer(s.nextID(), req.FirstName, req.LastName, req.Email, req.Phone, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.customerDAO.Create(ctx, customer); err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, customer.PullEvents()...); err != nil {
		return nil, err
	}

	return &CreateCustomerResp{ID: customer.GetID()}, nil
}
