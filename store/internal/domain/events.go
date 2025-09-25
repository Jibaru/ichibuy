package domain

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	StoreCreated    EventType = "StoreCreated"
	StoreUpdated    EventType = "StoreUpdated"
	StoreDeleted    EventType = "StoreDeleted"
	CustomerCreated EventType = "CustomerCreated"
	CustomerUpdated EventType = "CustomerUpdated"
	CustomerDeleted EventType = "CustomerDeleted"
	ProductCreated  EventType = "ProductCreated"
	ProductUpdated  EventType = "ProductUpdated"
	ProductDeleted  EventType = "ProductDeleted"
)

type Event struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

type EventBus interface {
	Publish(event Event) error
}

type StoreEventData struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Location    Location  `json:"location"`
	Slug        string    `json:"slug"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CustomerEventData struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductEventData struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Active      bool             `json:"active"`
	StoreID     string           `json:"store_id"`
	Images      map[string]Image `json:"images"`
	Prices      map[string]Price `json:"prices"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
