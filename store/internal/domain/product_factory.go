package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type ProductFactory struct {
	storageSvc StorageService
	nextID     NextID
}

func NewProductFactory(storageSvc StorageService, nextID NextID) *ProductFactory {
	return &ProductFactory{storageSvc: storageSvc, nextID: nextID}
}

func (f *ProductFactory) NewProduct(
	ctx context.Context,
	name string,
	description *string,
	active bool,
	storeID string,
	fileRequests []UploadFileRequest,
	prices []Price,
) (*Product, error) {
	// validations
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name is required")
	}

	if description != nil && strings.TrimSpace(*description) == "" {
		return nil, fmt.Errorf("description is required")
	}

	if len(fileRequests) == 0 {
		return nil, fmt.Errorf("at least one image is required")
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("at least one price is required")
	}

	uploadResp, err := f.storageSvc.UploadFiles(ctx, fileRequests)
	if err != nil {
		slog.ErrorContext(ctx, "upload files to storage failed", "error", err.Error())
		return nil, err
	}

	images := make([]Image, 0, len(uploadResp.Infos))
	for _, upload := range uploadResp.Infos {
		images = append(images, NewImage(upload.ID, upload.URL))
	}

	imagesMap := sliceToMap(images, func(v Image) string { return v.ID })
	pricesMap := sliceToMap(prices, func(v Price) string { return v.ID })

	rawImg, err := toRawMessage(imagesMap)
	if err != nil {
		return nil, err
	}

	rawPrice, err := toRawMessage(pricesMap)
	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:          f.nextID(),
		Name:        name,
		Description: description,
		Active:      active,
		StoreID:     storeID,
		Images:      rawImg,
		Prices:      rawPrice,
		images:      imagesMap,
		prices:      pricesMap,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	data, _ := json.Marshal(product.createEventData())
	event := Event{
		ID:        fmt.Sprintf("%s_%v", product.GetID(), product.CreatedAt.Unix()),
		Type:      ProductCreated,
		Data:      data,
		Timestamp: product.CreatedAt,
	}
	product.events = append(product.events, event)

	return product, nil
}
