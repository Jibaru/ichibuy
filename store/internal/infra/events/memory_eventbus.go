package events

import (
	"log"

	"ichibuy/store/internal/domain"
)

type MemoryEventBus struct{}

func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{}
}

func (e *MemoryEventBus) Publish(event domain.Event) error {
	log.Printf("Publishing event: %s with ID: %s at %s", event.Type, event.ID, event.Timestamp.Format("2006-01-02 15:04:05"))
	return nil
}