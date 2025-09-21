package dao

import (
	"context"
	"ichibuy/store/internal/domain"
)

type Store = domain.Store

type StoreDAO interface {
	// Create creates a new Store
	Create(ctx context.Context, m *Store) error

	// Update updates an existing Store
	Update(ctx context.Context, m *Store) error

	// PartialUpdate updates specific fields of a Store
	PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error

	// DeleteByPk deletes a Store by primary key
	DeleteByPk(ctx context.Context, pk string) error

	// FindByPk finds a Store by primary key
	FindByPk(ctx context.Context, pk string) (*Store, error)

	// CreateMany creates multiple Store records
	CreateMany(ctx context.Context, models []*Store) error

	// UpdateMany updates multiple Store records
	UpdateMany(ctx context.Context, models []*Store) error

	// DeleteManyByPks deletes multiple Store records by primary keys
	DeleteManyByPks(ctx context.Context, pks []string) error

	// FindOne finds a single Store with optional where clause and sort expression
	FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Store, error)

	// FindAll finds all Store records with optional where clause and sort expression
	FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Store, error)

	// FindPaginated finds Store records with pagination, optional where clause and sort expression
	FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Store, error)

	// Count counts Store records with optional where clause
	Count(ctx context.Context, where string, args ...interface{}) (int64, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
