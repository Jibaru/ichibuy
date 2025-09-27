package services

import (
	"context"
	"log/slog"

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
	slog.InfoContext(ctx, "update store started", "req", req)
	store, err := s.storeDAO.FindByPk(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "find store failed", "error", err.Error())
		return err
	}

	if err := store.Update(req.Name, req.Description, req.Location.Lat, req.Location.Lng, req.UserID); err != nil {
		slog.ErrorContext(ctx, "update store domain failed", "error", err.Error())
		return err
	}

	if err := s.storeDAO.Update(ctx, store); err != nil {
		slog.ErrorContext(ctx, "update store failed", "error", err.Error())
		return err
	}

	if err := s.eventBus.Publish(ctx, store.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return err
	}

	slog.InfoContext(ctx, "update store finished", "store_id", store.GetID())
	return nil
}
