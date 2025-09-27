package services

import (
	"context"
	"log/slog"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type CreateProductReq struct {
	Name        string
	Description *string
	Active      bool
	StoreID     string
	ImageFiles  []FileDTO
	Prices      []NewPriceDTO
}

type CreateProductResp = CreateUpdateResponse

type CreateProduct struct {
	productDAO     dao.ProductDAO
	eventBus       domain.EventBus
	nextID         domain.NextID
	productFactory *domain.ProductFactory
}

func NewCreateProduct(productDAO dao.ProductDAO, eventBus domain.EventBus, nextID domain.NextID, productFactory *domain.ProductFactory) *CreateProduct {
	return &CreateProduct{
		productDAO:     productDAO,
		eventBus:       eventBus,
		nextID:         nextID,
		productFactory: productFactory,
	}
}

func (s *CreateProduct) Exec(ctx context.Context, req CreateProductReq) (*CreateProductResp, error) {
	slog.InfoContext(ctx, "create product started", "req", req)

	prices, err := convertNewPriceDTOsToDomain(req.Prices, s.nextID)
	if err != nil {
		slog.ErrorContext(ctx, "convert new price dtos to domain failed", "error", err.Error())
		return nil, err
	}

	product, err := s.productFactory.NewProduct(ctx, req.Name, req.Description, req.Active, req.StoreID, s.fileDTOsToUploadFileRequests(req.ImageFiles), prices)
	if err != nil {
		slog.ErrorContext(ctx, "new product failed", "error", err.Error())
		return nil, err
	}

	if err := s.productDAO.Create(ctx, product); err != nil {
		slog.ErrorContext(ctx, "create product failed", "error", err.Error())
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, product.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return nil, err
	}
	slog.InfoContext(ctx, "create product finished", "product_id", product.GetID())
	return &CreateProductResp{ID: product.GetID()}, nil
}

func (s *CreateProduct) fileDTOsToUploadFileRequests(fileDTOs []FileDTO) []domain.UploadFileRequest {
	reqs := make([]domain.UploadFileRequest, len(fileDTOs))
	for i, fileDTO := range fileDTOs {
		reqs[i] = domain.UploadFileRequest{
			FileName:    fileDTO.FileName,
			ContentType: fileDTO.ContentType,
			Data:        fileDTO.Data,
		}
	}
	return reqs
}
