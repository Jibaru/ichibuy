package domain

import (
	"encoding/json"
	"fmt"
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

	Entity
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

func (p *Product) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Product) GetUpdatedAt() time.Time {
	return p.UpdatedAt
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

	data, _ := json.Marshal(p.createEventData())
	event := Event{
		ID:        fmt.Sprintf("%s_%v", p.GetID(), p.GetUpdatedAt().Unix()),
		Type:      ProductUpdated,
		Data:      data,
		Timestamp: p.GetUpdatedAt(),
	}
	p.events = append(p.events, event)

	return nil
}

func (p *Product) PrepareDelete() {
	data, _ := json.Marshal(p.createEventData())

	event := Event{
		ID:        fmt.Sprintf("%s_%v_delete", p.GetID(), p.GetUpdatedAt().Unix()),
		Type:      ProductDeleted,
		Data:      data,
		Timestamp: p.GetUpdatedAt(),
	}

	p.events = append(p.events, event)
}

func (p *Product) createEventData() ProductEventData {
	return ProductEventData{
		ID:          p.GetID(),
		Name:        p.GetName(),
		Description: p.GetDescription(),
		Active:      p.GetActive(),
		StoreID:     p.GetStoreID(),
		Images:      p.GetImages(),
		Prices:      p.GetPrices(),
		CreatedAt:   p.GetCreatedAt(),
		UpdatedAt:   p.GetUpdatedAt(),
	}
}
