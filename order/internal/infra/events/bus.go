package events

import (
	"context"

	"ichibuy/order/internal/domain"
	"ichibuy/order/internal/domain/dao"
)

type Bus struct {
	eventDAO dao.EventDAO
}

func NewBus(eventDAO dao.EventDAO) *Bus {
	return &Bus{
		eventDAO: eventDAO,
	}
}

func (b *Bus) Publish(ctx context.Context, events ...domain.Event) error {
	evts := make([]*domain.Event, len(events))
	for i, event := range events {
		evts[i] = &event
	}

	return b.eventDAO.CreateMany(ctx, evts)
}
