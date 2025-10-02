package services

import (
	"context"
	"log/slog"
	"time"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type GetCustomerByUserIDReq struct {
	UserID string
}

type GetCustomerByUserIDResp struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetCustomerByUserID struct {
	customerDAO dao.CustomerDAO
}

func NewGetCustomerByUserID(customerDAO dao.CustomerDAO) *GetCustomerByUserID {
	return &GetCustomerByUserID{
		customerDAO: customerDAO,
	}
}

func (s *GetCustomerByUserID) Exec(ctx context.Context, req GetCustomerByUserIDReq) (*GetCustomerByUserIDResp, error) {
	slog.InfoContext(ctx, "get customer by user id started", "req", req)
	customer, err := s.customerDAO.FindOne(ctx, "user_id = $1", "", req.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "find customer by user id failed", "error", err.Error())
		return nil, err
	}
	slog.InfoContext(ctx, "get customer by user id finished", "customer_id", customer.GetID())
	return mapCustomerToGetCustomerByUserIDResp(customer), nil
}

func mapCustomerToGetCustomerByUserIDResp(customer *domain.Customer) *GetCustomerByUserIDResp {
	return &GetCustomerByUserIDResp{
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
