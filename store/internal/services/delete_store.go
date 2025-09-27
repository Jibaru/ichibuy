package services

import (
	"context"

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

	store.PrepareDelete()

	if err := s.storeDAO.DeleteByPk(ctx, req.ID); err != nil {
		return err
	}

	if err := s.eventBus.Publish(ctx, store.PullEvents()...); err != nil {
		return err
	}

	return nil
}
