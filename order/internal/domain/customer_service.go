package domain

import "context"

type CustomerService interface {
	FindByUserID(ctx context.Context, userID string) (*CustomerDTO, error)
}

type CustomerDTO struct {
	ID string
}
