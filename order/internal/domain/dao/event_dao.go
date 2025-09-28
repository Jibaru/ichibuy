package dao

import (
	"context"
	"ichibuy/order/internal/domain"
)

type Event = domain.Event

type EventDAO interface {
	// Create creates a new Event
	Create(ctx context.Context, m *Event) error

	// Update updates an existing Event
	Update(ctx context.Context, m *Event) error

	// PartialUpdate updates specific fields of a Event
	PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error

	// DeleteByPk deletes a Event by primary key
	DeleteByPk(ctx context.Context, pk string) error

	// FindByPk finds a Event by primary key
	FindByPk(ctx context.Context, pk string) (*Event, error)

	// CreateMany creates multiple Event records
	CreateMany(ctx context.Context, models []*Event) error

	// UpdateMany updates multiple Event records
	UpdateMany(ctx context.Context, models []*Event) error

	// DeleteManyByPks deletes multiple Event records by primary keys
	DeleteManyByPks(ctx context.Context, pks []string) error

	// FindOne finds a single Event with optional where clause and sort expression
	FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Event, error)

	// FindAll finds all Event records with optional where clause and sort expression
	FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Event, error)

	// FindPaginated finds Event records with pagination, optional where clause and sort expression
	FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Event, error)

	// Count counts Event records with optional where clause
	Count(ctx context.Context, where string, args ...interface{}) (int64, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
