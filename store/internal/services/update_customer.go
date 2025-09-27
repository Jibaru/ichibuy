package services

import (
	"context"

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
	customer, err := s.customerDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return err
	}

	if err := customer.Update(req.FirstName, req.LastName, req.Email, req.Phone, req.UserID); err != nil {
		return err
	}

	if err := s.customerDAO.Update(ctx, customer); err != nil {
		return err
	}

	if err := s.eventBus.Publish(ctx, customer.PullEvents()...); err != nil {
		return err
	}

	return nil
}
