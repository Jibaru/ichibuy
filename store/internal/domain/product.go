package domain

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID          string          `sql:"id,primary"`
	Name        string          `sql:"name"`
	Description *string         `sql:"description"`
	Active      bool            `sql:"active"`
	StoreID     string          `sql:"store_id"`
	Images      json.RawMessage `sql:"images"`
	Prices      json.RawMessage `sql:"prices"`
	CreatedAt   time.Time       `sql:"created_at"`
	UpdatedAt   time.Time       `sql:"updated_at"`

	// Non-storable
	prices map[string]Price
	images map[string]Image
}

type Image struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Price struct {
	ID    string `json:"id"`
	Value Money  `json:"value"`
}

func sliceToMap[T any](slice []T, id func(v T) string) map[string]T {
	m := make(map[string]T)
	for _, item := range slice {
		m[id(item)] = item
	}
	return m
}

func NewProduct(id, name string, description *string, active bool, storeID string, images []Image, prices []Price) (*Product, error) {
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

	return &Product{
		ID:          id,
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
	}, nil
}

func NewImage(id, url string) Image {
	return Image{
		ID:  id,
		URL: url,
	}
}

func NewPrice(id string, value Money) Price {
	return Price{
		ID:    id,
		Value: value,
	}
}

func toRawMessage(v interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) GetID() string {
	return p.ID
}

func (p *Product) GetName() string {
	return p.Name
}

func (p *Product) GetDescription() *string {
	return p.Description
}

func (p *Product) GetActive() bool {
	return p.Active
}

func (p *Product) GetStoreID() string {
	return p.StoreID
}

func (p *Product) GetImages() map[string]Image {
	if p.images == nil {
		_ = json.Unmarshal(p.Images, &p.images)
	}

	return p.images
}

func (p *Product) GetPrices() map[string]Price {
	if p.prices == nil {
		_ = json.Unmarshal(p.Prices, &p.prices)
	}

	return p.prices
}

func (p *Product) Update(
	name string,
	description *string,
	active bool,
	newImages []Image,
	newPrices []Price,
	deleteImagesIDs []string,
	deletePricesIDs []string,
) error {
	p.Name = name
	p.Description = description
	p.Active = active
	p.UpdatedAt = time.Now().UTC()

	for _, img := range newImages {
		p.images[img.ID] = img
	}

	for _, id := range deleteImagesIDs {
		delete(p.images, id)
	}

	for _, price := range newPrices {
		p.prices[price.ID] = price
	}

	for _, id := range deletePricesIDs {
		delete(p.prices, id)
	}

	rawImg, err := toRawMessage(p.images)
	if err != nil {
		return err
	}

	rawPrice, err := toRawMessage(p.prices)
	if err != nil {
		return err
	}

	p.Images = rawImg
	p.Prices = rawPrice

	return nil
}
