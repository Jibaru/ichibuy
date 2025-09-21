package services

import (
	"context"
	"time"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type GetCustomerReq struct {
	ID string
}

type GetCustomerResp struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetCustomer struct {
	customerDAO dao.CustomerDAO
}

func NewGetCustomer(customerDAO dao.CustomerDAO) *GetCustomer {
	return &GetCustomer{
		customerDAO: customerDAO,
	}
}

func (s *GetCustomer) Exec(ctx context.Context, req GetCustomerReq) (*GetCustomerResp, error) {
	customer, err := s.customerDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapCustomerToGetCustomerResp(customer), nil
}

func mapCustomerToGetCustomerResp(customer *domain.Customer) *GetCustomerResp {
	return &GetCustomerResp{
		ID:        customer.GetID(),
		FirstName: customer.GetFirstName(),
		LastName:  customer.GetLastName(),
		Email:     customer.GetEmailString(),
		Phone:     customer.GetPhoneString(),
		UserID:    customer.GetUserID(),
		CreatedAt: customer.GetCreatedAt(),
		UpdatedAt: customer.GetUpdatedAt(),
	}
}
