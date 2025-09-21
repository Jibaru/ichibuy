package services

import (
	"context"
	"time"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type GetStoreReq struct {
	ID string
}

type GetStoreResp struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Slug        string    `json:"slug"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetStore struct {
	storeDAO dao.StoreDAO
}

func NewGetStore(storeDAO dao.StoreDAO) *GetStore {
	return &GetStore{
		storeDAO: storeDAO,
	}
}

func (s *GetStore) Exec(ctx context.Context, req GetStoreReq) (*GetStoreResp, error) {
	store, err := s.storeDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapStoreToGetStoreResp(store), nil
}

func mapStoreToGetStoreResp(store *domain.Store) *GetStoreResp {
	return &GetStoreResp{
		ID:          store.GetID(),
		Name:        store.GetName(),
		Description: store.GetDescription(),
		Lat:         store.GetLat(),
		Lng:         store.GetLng(),
		Slug:        store.GetSlug(),
		UserID:      store.GetUserID(),
		CreatedAt:   store.GetCreatedAt(),
		UpdatedAt:   store.GetUpdatedAt(),
	}
}
