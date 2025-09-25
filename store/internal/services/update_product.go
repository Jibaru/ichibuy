package services

import (
	"context"
	"encoding/json"

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
	product, err := s.productDAO.FindByPk(ctx, req.ID)
	if err != nil {
		return err
	}

	// Upload new images to storage
	images, err := s.uploadImages(ctx, req.NewImageFiles)
	if err != nil {
		return err
	}

	if err := s.storageSvc.DeleteFiles(ctx, req.DeleteImageIDs); err != nil {
		return err
	}

	prices, err := convertNewPriceDTOsToDomain(req.NewPrices, s.nextID)
	if err != nil {
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
		return err
	}

	if err := s.productDAO.Update(ctx, product); err != nil {
		return err
	}

	eventData := domain.ProductEventData{
		ID:          product.GetID(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Active:      product.GetActive(),
		StoreID:     product.GetStoreID(),
		Images:      product.GetImages(),
		Prices:      product.GetPrices(),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	data, _ := json.Marshal(eventData)
	event := domain.Event{
		ID:        s.nextID(),
		Type:      domain.ProductUpdated,
		Data:      data,
		Timestamp: product.UpdatedAt,
	}

	s.eventBus.Publish(event)

	return nil
}

func (s *UpdateProduct) uploadImages(ctx context.Context, fileDTOs []FileDTO) ([]domain.Image, error) {
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
		return nil, err
	}

	images := make([]domain.Image, 0, len(fileDTOs))
	for _, u := range uploadResp.Infos {
		images = append(images, domain.NewImage(u.ID, u.URL))
	}

	return images, nil
}

func convertNewPriceDTOsToDomain(dtos []NewPriceDTO, nextID domain.NextID) ([]domain.Price, error) {
	prices := make([]domain.Price, len(dtos))
	for i, dto := range dtos {
		money, err := domain.NewMoney(dto.Amount, dto.Currency)
		if err != nil {
			return nil, err
		}
		prices[i] = domain.NewPrice(nextID(), money)
	}
	return prices, nil
}
