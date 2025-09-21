package services

import (
	"context"
	"encoding/json"

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
		Type:      domain.CustomerUpdated,
		Data:      data,
		Timestamp: customer.GetUpdatedAt(),
	}

	s.eventBus.Publish(event)

	return nil
}
