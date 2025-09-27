package services

import (
	"context"
	"log/slog"

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
	slog.InfoContext(ctx, "delete store started", "req", req)
	store, err := s.storeDAO.FindByPk(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "find store failed", "error", err.Error())
		return err
	}

	store.PrepareDelete()

	if err := s.storeDAO.DeleteByPk(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "delete store failed", "error", err.Error())
		return err
	}

	if err := s.eventBus.Publish(ctx, store.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return err
	}

	slog.InfoContext(ctx, "delete store finished", "store_id", req.ID)
	return nil
}
