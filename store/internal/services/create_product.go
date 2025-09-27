package services

import (
	"context"

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
	images, err := s.uploadImages(ctx, req.ImageFiles)
	if err != nil {
		return nil, err
	}

	prices, err := convertNewPriceDTOsToDomain(req.Prices, s.nextID)
	if err != nil {
		return nil, err
	}

	product, err := domain.NewProduct(s.nextID(), req.Name, req.Description, req.Active, req.StoreID, images, prices)
	if err != nil {
		return nil, err
	}

	if err := s.productDAO.Create(ctx, product); err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, product.PullEvents()...); err != nil {
		return nil, err
	}

	return &CreateProductResp{ID: product.GetID()}, nil
}

func (s *CreateProduct) uploadImages(ctx context.Context, fileDTOs []FileDTO) ([]domain.Image, error) {
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
	for _, upload := range uploadResp.Infos {
		images = append(images, domain.NewImage(upload.ID, upload.URL))
	}

	return images, nil
}
