package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type ListProductsReq struct {
	Filters    ProductFilters
	Pagination Pagination
	Sorting    Sorting
}

type ProductFilters struct {
	StoreID     string
	Name        *string
	Description *string
	Active      *bool
}

type ProductListItem struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Active      bool       `json:"active"`
	StoreID     string     `json:"store_id"`
	Images      []ImageDTO `json:"images"`
	Prices      []PriceDTO `json:"prices"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ListProductsResp struct {
	Products []ProductListItem `json:"products"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type ListProducts struct {
	productDAO dao.ProductDAO
}

func NewListProducts(productDAO dao.ProductDAO) *ListProducts {
	return &ListProducts{
		productDAO: productDAO,
	}
}

func (s *ListProducts) Exec(ctx context.Context, req ListProductsReq) (ListProductsResp, error) {
	var whereParts []string
	var args []any
	i := 1

	if req.Filters.StoreID != "" {
		whereParts = append(whereParts, fmt.Sprintf("store_id = $%d", i))
		args = append(args, req.Filters.StoreID)
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

	if req.Filters.Active != nil {
		whereParts = append(whereParts, fmt.Sprintf("active = $%d", i))
		args = append(args, *req.Filters.Active)
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

	total, err := s.productDAO.Count(ctx, where, args...)
	if err != nil {
		return ListProductsResp{}, err
	}

	products, err := s.productDAO.FindPaginated(
		ctx,
		req.Pagination.Limit,
		req.Pagination.Offset,
		where,
		sort,
		args...,
	)
	if err != nil {
		return ListProductsResp{}, err
	}

	return ListProductsResp{
		Products: mapProductsToListProductsResp(products),
		Total:    total,
		Limit:    req.Pagination.Limit,
		Offset:   req.Pagination.Offset,
	}, nil
}

func mapProductsToListProductsResp(products []*domain.Product) []ProductListItem {
	response := make([]ProductListItem, len(products))
	for i, product := range products {
		response[i] = ProductListItem{
			ID:          product.GetID(),
			Name:        product.GetName(),
			Description: product.GetDescription(),
			Active:      product.GetActive(),
			StoreID:     product.GetStoreID(),
			Images:      convertDomainImagesToDTOs(product.GetImages()),
			Prices:      convertDomainPricesToDTOs(product.GetPrices()),
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}
	return response
}

func convertDomainImagesToDTOs(images map[string]domain.Image) []ImageDTO {
	dtos := make([]ImageDTO, len(images))
	i := 0
	for _, img := range images {
		dtos[i] = ImageDTO{
			ID:  img.ID,
			URL: img.URL,
		}
		i++
	}
	return dtos
}

func convertDomainPricesToDTOs(prices map[string]domain.Price) []PriceDTO {
	dtos := make([]PriceDTO, len(prices))
	i := 0
	for _, price := range prices {
		dtos[i] = PriceDTO{
			ID:       price.ID,
			Amount:   price.Value.GetAmount(),
			Currency: price.Value.GetCurrency(),
		}
		i++
	}
	return dtos
}
