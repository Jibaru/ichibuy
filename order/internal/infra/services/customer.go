package services

import (
	"context"

	storeHTTP "github.com/Jibaru/ichibuy/api-client/go/store"

	"ichibuy/order/internal/domain"
	sharedCtx "ichibuy/order/internal/shared/context"
)

type customerService struct {
	client *storeHTTP.APIClient
}

func NewCustomerService(client *storeHTTP.APIClient) *customerService {
	return &customerService{client: client}
}

func (s *customerService) FindByUserID(ctx context.Context, userID string) (*domain.CustomerDTO, error) {
	ctx = sharedCtx.AddToken(ctx, storeHTTP.ContextAccessToken)

	resp, _, err := s.client.CustomersApi.ApiV1CustomersUserUserIdGet(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.CustomerDTO{
		ID: resp.Id,
	}, nil
}
