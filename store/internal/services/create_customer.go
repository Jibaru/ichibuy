package services

import (
	"context"
	"encoding/json"

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
		Type:      domain.CustomerCreated,
		Data:      data,
		Timestamp: customer.GetCreatedAt(),
	}

	s.eventBus.Publish(event)

	return &CreateCustomerResp{ID: customer.GetID()}, nil
}
