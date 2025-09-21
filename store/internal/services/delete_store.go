package services

import (
	"context"
	"encoding/json"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type DeleteStoreReq struct {
	ID string
}

type DeleteStore struct {
	storeDAO dao.StoreDAO
	eventBus domain.EventBus
	nextID   domain.NextID
}

func NewDeleteStore(storeDAO dao.StoreDAO, eventBus domain.EventBus, nextID domain.NextID) *DeleteStore {
	return &DeleteStore{
		storeDAO: storeDAO,
		eventBus: eventBus,
		nextID:   nextID,
	}
}

func (s *DeleteStore) Exec(ctx context.Context, req DeleteStoreReq) error {
	store, err := s.storeDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return err
	}

	if err := s.storeDAO.DeleteByPk(ctx, req.ID); err != nil {
		return err
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
		Type:      domain.StoreDeleted,
		Data:      data,
		Timestamp: store.GetUpdatedAt(),
	}

	s.eventBus.Publish(event)

	return nil
}
