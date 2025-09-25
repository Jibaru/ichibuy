package dao

import (
	"context"
	"ichibuy/store/internal/domain"
)

type Product = domain.Product

type ProductDAO interface {
	// Create creates a new Product
	Create(ctx context.Context, m *Product) error

	// Update updates an existing Product
	Update(ctx context.Context, m *Product) error

	// PartialUpdate updates specific fields of a Product
	PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error

	// DeleteByPk deletes a Product by primary key
	DeleteByPk(ctx context.Context, pk string) error

	// FindByPk finds a Product by primary key
	FindByPk(ctx context.Context, pk string) (*Product, error)

	// CreateMany creates multiple Product records
	CreateMany(ctx context.Context, models []*Product) error

	// UpdateMany updates multiple Product records
	UpdateMany(ctx context.Context, models []*Product) error

	// DeleteManyByPks deletes multiple Product records by primary keys
	DeleteManyByPks(ctx context.Context, pks []string) error

	// FindOne finds a single Product with optional where clause and sort expression
	FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Product, error)

	// FindAll finds all Product records with optional where clause and sort expression
	FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Product, error)

	// FindPaginated finds Product records with pagination, optional where clause and sort expression
	FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Product, error)

	// Count counts Product records with optional where clause
	Count(ctx context.Context, where string, args ...interface{}) (int64, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
