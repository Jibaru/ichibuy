package services

import (
	"context"
	"log/slog"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/domain/dao"
)

type UpdateProductReq struct {
	ID              string
	Name            string
	Description     *string
	Active          bool
	NewImageFiles   []FileDTO
	DeleteImageIDs  []string
	NewPrices       []NewPriceDTO
	DeletePricesIDs []string
}

type UpdateProduct struct {
	productDAO dao.ProductDAO
	eventBus   domain.EventBus
	nextID     domain.NextID
	storageSvc domain.StorageService
}

func NewUpdateProduct(productDAO dao.ProductDAO, eventBus domain.EventBus, nextID domain.NextID, storageSvc domain.StorageService) *UpdateProduct {
	return &UpdateProduct{
		productDAO: productDAO,
		eventBus:   eventBus,
		nextID:     nextID,
		storageSvc: storageSvc,
	}
}

func (s *UpdateProduct) Exec(ctx context.Context, req UpdateProductReq) error {
	slog.InfoContext(ctx, "update product started", "req", req)
	product, err := s.productDAO.FindByPk(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "find product failed", "error", err.Error())
		return err
	}

	// Upload new images to storage
	images, err := s.uploadImages(ctx, req.NewImageFiles)
	if err != nil {
		slog.ErrorContext(ctx, "upload images failed", "error", err.Error())
		return err
	}

	if err := s.storageSvc.DeleteFiles(ctx, req.DeleteImageIDs); err != nil {
		slog.ErrorContext(ctx, "delete files from storage failed", "error", err.Error())
		return err
	}

	prices, err := convertNewPriceDTOsToDomain(req.NewPrices, s.nextID)
	if err != nil {
		slog.ErrorContext(ctx, "convert new price dtos to domain failed", "error", err.Error())
		return err
	}

	err = product.Update(
		req.Name,
		req.Description,
		req.Active,
		images,
		prices,
		req.DeleteImageIDs,
		req.DeletePricesIDs,
	)
	if err != nil {
		slog.ErrorContext(ctx, "update product domain failed", "error", err.Error())
		return err
	}

	if err := s.productDAO.Update(ctx, product); err != nil {
		slog.ErrorContext(ctx, "update product failed", "error", err.Error())
		return err
	}

	if err := s.eventBus.Publish(ctx, product.PullEvents()...); err != nil {
		slog.ErrorContext(ctx, "publish events failed", "error", err.Error())
		return err
	}

	slog.InfoContext(ctx, "update product finished", "product_id", product.GetID())
	return nil
}

func (s *UpdateProduct) uploadImages(ctx context.Context, fileDTOs []FileDTO) ([]domain.Image, error) {
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
	for _, u := range uploadResp.Infos {
		images = append(images, domain.NewImage(u.ID, u.URL))
	}

	slog.InfoContext(ctx, "images uploaded", "image_count", len(images))
	return images, nil
}

func convertNewPriceDTOsToDomain(dtos []NewPriceDTO, nextID domain.NextID) ([]domain.Price, error) {
	prices := make([]domain.Price, len(dtos))
	for i, dto := range dtos {
		money, err := domain.NewMoney(dto.Amount, dto.Currency)
		if err != nil {
			slog.Error("new money failed", "error", err.Error())
			return nil, err
		}
		prices[i] = domain.NewPrice(nextID(), money)
	}
	return prices, nil
}
