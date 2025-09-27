package services

import (
	"context"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type CreateStoreReq struct {
	Name        string
	Description *string
	Location    domain.Location
	UserID      string
}

type CreateStoreResp = CreateUpdateResponse

type CreateStore struct {
	storeDAO dao.StoreDAO
	eventBus domain.EventBus
	nextID   domain.NextID
}

func NewCreateStore(storeDAO dao.StoreDAO, eventBus domain.EventBus, nextID domain.NextID) *CreateStore {
	return &CreateStore{
		storeDAO: storeDAO,
		eventBus: eventBus,
		nextID:   nextID,
	}
}

func (s *CreateStore) Exec(ctx context.Context, req CreateStoreReq) (*CreateStoreResp, error) {
	store, err := domain.NewStore(s.nextID(), req.Name, req.Description, req.Location.Lat, req.Location.Lng, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.storeDAO.Create(ctx, store); err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, store.PullEvents()...); err != nil {
		return nil, err
	}

	return &CreateStoreResp{ID: store.GetID()}, nil
}
