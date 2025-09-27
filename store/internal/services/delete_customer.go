package services

import (
	"context"

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
	customer, err := s.customerDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return err
	}

	customer.PrepareDelete()

	if err := s.customerDAO.DeleteByPk(ctx, req.ID); err != nil {
		return err
	}

	if err := s.eventBus.Publish(ctx, customer.PullEvents()...); err != nil {
		return err
	}

	return nil
}
