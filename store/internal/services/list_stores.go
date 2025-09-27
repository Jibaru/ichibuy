package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type ListStoresReq struct {
	Filters    StoreFilters
	Pagination Pagination
	Sorting    Sorting
}

type StoreFilters struct {
	UserID      string
	Name        *string
	Description *string
}

type StoreListItem struct {
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

type ListStoresResp struct {
	Stores []StoreListItem `json:"stores"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

type ListStores struct {
	storeDAO dao.StoreDAO
}

func NewListStores(storeDAO dao.StoreDAO) *ListStores {
	return &ListStores{
		storeDAO: storeDAO,
	}
}

func (s *ListStores) Exec(ctx context.Context, req ListStoresReq) (ListStoresResp, error) {
	slog.InfoContext(ctx, "list stores started", "req", req)
	var whereParts []string
	var args []any
	i := 1

	if req.Filters.UserID != "" {
		whereParts = append(whereParts, fmt.Sprintf("user_id = $%d", i))
		args = append(args, req.Filters.UserID)
		i++
	}

	if req.Filters.Name != nil {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", i))
		args = append(args, "%"+*req.Filters.Name+"%")
		i++
	}

	if req.Filters.Description != nil {
		whereParts = append(whereParts, fmt.Sprintf("description ILIKE $%d", i))
		args = append(args, "%"+*req.Filters.Description+"%")
		i++
	}

	where := strings.Join(whereParts, " AND ")

	sort := ""
	if req.Sorting.Field != "" {
		order := "ASC"
		if strings.ToUpper(req.Sorting.Order) == "DESC" {
			order = "DESC"
		}
		sort = fmt.Sprintf("%s %s", req.Sorting.Field, order)
	}

	total, err := s.storeDAO.Count(ctx, where, args...)
	if err != nil {
		slog.ErrorContext(ctx, "count stores failed", "error", err.Error())
		return ListStoresResp{}, err
	}

	stores, err := s.storeDAO.FindPaginated(
		ctx,
		req.Pagination.Limit,
		req.Pagination.Offset,
		where,
		sort,
		args...,
	)
	if err != nil {
		slog.ErrorContext(ctx, "find paginated stores failed", "error", err.Error())
		return ListStoresResp{}, err
	}

	slog.InfoContext(ctx, "list stores finished", "total", total, "count", len(stores))
	return ListStoresResp{
		Stores: mapStoresToListStoresResp(stores),
		Total:  total,
		Limit:  req.Pagination.Limit,
		Offset: req.Pagination.Offset,
	}, nil
}

func mapStoresToListStoresResp(stores []*domain.Store) []StoreListItem {
	response := make([]StoreListItem, len(stores))
	for i, store := range stores {
		response[i] = StoreListItem{
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
	return response
}
