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
	productDAO dao.ProductDAO
	eventBus   domain.EventBus
	nextID     domain.NextID
	storageSvc domain.StorageService
}

func NewCreateProduct(productDAO dao.ProductDAO, eventBus domain.EventBus, nextID domain.NextID, storageSvc domain.StorageService) *CreateProduct {
	return &CreateProduct{
		productDAO: productDAO,
		eventBus:   eventBus,
		nextID:     nextID,
		storageSvc: storageSvc,
	}
}

func (s *CreateProduct) Exec(ctx context.Context, req CreateProductReq) (*CreateProductResp, error) {
	slog.InfoContext(ctx, "create product started", "req", req)
	images, err := s.uploadImages(ctx, req.ImageFiles)
	if err != nil {
		slog.ErrorContext(ctx, "upload images failed", "error", err.Error())
		return nil, err
	}

	prices, err := convertNewPriceDTOsToDomain(req.Prices, s.nextID)
	if err != nil {
		slog.ErrorContext(ctx, "convert new price dtos to domain failed", "error", err.Error())
		return nil, err
	}

	product, err := domain.NewProduct(s.nextID(), req.Name, req.Description, req.Active, req.StoreID, images, prices)
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

func (s *CreateProduct) uploadImages(ctx context.Context, fileDTOs []FileDTO) ([]domain.Image, error) {
	slog.InfoContext(ctx, "uploading images", "file_count", len(fileDTOs))
	reqs := make([]domain.UploadFileRequest, len(fileDTOs))
	for i, fileDTO := range fileDTOs {
		reqs[i] = domain.UploadFileRequest{
			FileName:    fileDTO.FileName,
			ContentType: fileDTO.ContentType,
			Data:        fileDTO.Data,
		}
	}

	uploadResp, err := s.storageSvc.UploadFiles(ctx, reqs)
	if err != nil {
		slog.ErrorContext(ctx, "upload files to storage failed", "error", err.Error())
		return nil, err
	}

	images := make([]domain.Image, 0, len(fileDTOs))
	for _, upload := range uploadResp.Infos {
		images = append(images, domain.NewImage(upload.ID, upload.URL))
	}
	slog.InfoContext(ctx, "images uploaded", "image_count", len(images))
	return images, nil
}
