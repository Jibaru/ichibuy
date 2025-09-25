package services

import (
	"context"

	"ichibuy/store/internal/domain/dao"
)

type GetProductResp struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Active      bool       `json:"active"`
	StoreID     string     `json:"store_id"`
	Images      []ImageDTO `json:"images"`
	Prices      []PriceDTO `json:"prices"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

type GetProduct struct {
	productDAO dao.ProductDAO
}

func NewGetProduct(productDAO dao.ProductDAO) *GetProduct {
	return &GetProduct{
		productDAO: productDAO,
	}
}

func (s *GetProduct) Exec(ctx context.Context, id string) (*GetProductResp, error) {
	product, err := s.productDAO.FindByPk(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetProductResp{
		ID:          product.GetID(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Active:      product.GetActive(),
		StoreID:     product.GetStoreID(),
		Images:      convertDomainImagesToDTOs(product.GetImages()),
		Prices:      convertDomainPricesToDTOs(product.GetPrices()),
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
