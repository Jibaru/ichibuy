package domain

import (
	"context"
	"encoding/json"
	"time"
)

type EventType string

const (
	OrderCreated EventType = "OrderCreated"
)

type Event struct {
	ID        string          `json:"id" sql:"id,primary"`
	Type      EventType       `json:"type" sql:"type"`
	Data      json.RawMessage `json:"data" sql:"data"`
	Timestamp time.Time       `json:"timestamp" sql:"timestamp"`
}

func (e *Event) TableName() string {
	return "events"
}

type EventBus interface {
	Publish(ctx context.Context, events ...Event) error
}
