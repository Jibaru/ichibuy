package services

import (
	"context"
	"encoding/json"

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

	eventData := domain.StoreEventData{
		ID:          store.GetID(),
		Name:        store.GetName(),
		Description: store.GetDescription(),
		Location:    store.Location(),
		Slug:        store.GetSlug(),
		UserID:      store.GetUserID(),
		CreatedAt:   store.GetCreatedAt(),
		UpdatedAt:   store.GetUpdatedAt(),
	}

	data, _ := json.Marshal(eventData)
	event := domain.Event{
		ID:        s.nextID(),
		Type:      domain.StoreCreated,
		Data:      data,
		Timestamp: store.GetCreatedAt(),
	}

	s.eventBus.Publish(event)

	return &CreateStoreResp{ID: store.GetID()}, nil
}
