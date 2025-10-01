package dao

import (
	"context"
	"ichibuy/order/internal/domain"
)

type Order = domain.Order

type OrderDAO interface {
	// Create creates a new Order
	Create(ctx context.Context, m *Order) error

	// Update updates an existing Order
	Update(ctx context.Context, m *Order) error

	// PartialUpdate updates specific fields of a Order
	PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error

	// DeleteByPk deletes a Order by primary key
	DeleteByPk(ctx context.Context, pk string) error

	// FindByPk finds a Order by primary key
	FindByPk(ctx context.Context, pk string) (*Order, error)

	// CreateMany creates multiple Order records
	CreateMany(ctx context.Context, models []*Order) error

	// UpdateMany updates multiple Order records
	UpdateMany(ctx context.Context, models []*Order) error

	// DeleteManyByPks deletes multiple Order records by primary keys
	DeleteManyByPks(ctx context.Context, pks []string) error

	// FindOne finds a single Order with optional where clause and sort expression
	FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Order, error)

	// FindAll finds all Order records with optional where clause and sort expression
	FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Order, error)

	// FindPaginated finds Order records with pagination, optional where clause and sort expression
	FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Order, error)

	// Count counts Order records with optional where clause
	Count(ctx context.Context, where string, args ...interface{}) (int64, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
