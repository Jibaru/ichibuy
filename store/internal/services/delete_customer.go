package services

import (
	"context"
	"encoding/json"

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

	if err := s.customerDAO.DeleteByPk(ctx, req.ID); err != nil {
		return err
	}

	var emailStr, phoneStr *string
	if customer.GetEmail() != nil {
		emailValue := customer.GetEmail().Value()
		emailStr = &emailValue
	}
	if customer.GetPhone() != nil {
		phoneValue := customer.GetPhone().Value()
		phoneStr = &phoneValue
	}

	eventData := domain.CustomerEventData{
		ID:        customer.GetID(),
		FirstName: customer.GetFirstName(),
		LastName:  customer.GetLastName(),
		Email:     emailStr,
		Phone:     phoneStr,
		UserID:    customer.GetUserID(),
		CreatedAt: customer.GetCreatedAt(),
		UpdatedAt: customer.GetUpdatedAt(),
	}

	data, _ := json.Marshal(eventData)
	event := domain.Event{
		ID:        s.nextID(),
		Type:      domain.CustomerDeleted,
		Data:      data,
		Timestamp: customer.GetUpdatedAt(),
	}

	s.eventBus.Publish(event)

	return nil
}
