package services

import (
	"context"
	"encoding/json"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type UpdateStoreReq struct {
	ID          string
	Name        string
	Description *string
	Location    domain.Location
	UserID      string
}

type UpdateStore struct {
	storeDAO dao.StoreDAO
	eventBus domain.EventBus
	nextID   domain.NextID
}

func NewUpdateStore(storeDAO dao.StoreDAO, eventBus domain.EventBus, nextID domain.NextID) *UpdateStore {
	return &UpdateStore{
		storeDAO: storeDAO,
		eventBus: eventBus,
		nextID:   nextID,
	}
}

func (s *UpdateStore) Exec(ctx context.Context, req UpdateStoreReq) error {
	store, err := s.storeDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return err
	}

	if err := store.Update(req.Name, req.Description, req.Location.Lat, req.Location.Lng, req.UserID); err != nil {
		return err
	}

	if err := s.storeDAO.Update(ctx, store); err != nil {
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
		Type:      domain.StoreUpdated,
		Data:      data,
		Timestamp: store.GetUpdatedAt(),
	}

	s.eventBus.Publish(event)

	return nil
}
