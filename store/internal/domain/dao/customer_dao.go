package dao

import (
	"context"
	"ichibuy/store/internal/domain"
)

type Customer = domain.Customer

type CustomerDAO interface {
	// Create creates a new Customer
	Create(ctx context.Context, m *Customer) error

	// Update updates an existing Customer
	Update(ctx context.Context, m *Customer) error

	// PartialUpdate updates specific fields of a Customer
	PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error

	// DeleteByPk deletes a Customer by primary key
	DeleteByPk(ctx context.Context, pk string) error

	// FindByPk finds a Customer by primary key
	FindByPk(ctx context.Context, pk string) (*Customer, error)

	// CreateMany creates multiple Customer records
	CreateMany(ctx context.Context, models []*Customer) error

	// UpdateMany updates multiple Customer records
	UpdateMany(ctx context.Context, models []*Customer) error

	// DeleteManyByPks deletes multiple Customer records by primary keys
	DeleteManyByPks(ctx context.Context, pks []string) error

	// FindOne finds a single Customer with optional where clause and sort expression
	FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Customer, error)

	// FindAll finds all Customer records with optional where clause and sort expression
	FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Customer, error)

	// FindPaginated finds Customer records with pagination, optional where clause and sort expression
	FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Customer, error)

	// Count counts Customer records with optional where clause
	Count(ctx context.Context, where string, args ...interface{}) (int64, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
