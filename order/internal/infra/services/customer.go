package services

import (
	"context"

	storeHTTP "github.com/Jibaru/ichibuy/api-client/go/store"

	"ichibuy/order/internal/domain"
)

type customerService struct {
	client *storeHTTP.APIClient
}

func NewCustomerService(client *storeHTTP.APIClient) *customerService {
	return &customerService{client: client}
}

func (s *customerService) FindByUserID(ctx context.Context, userID string) (*domain.CustomerDTO, error) {
	// TODO: implement
	return &domain.CustomerDTO{}, nil
}
