package dao

import (
	"context"
	"ichibuy/auth/internal/domain"
)

type User = domain.User

type UserDAO interface {
	// Create creates a new User
	Create(ctx context.Context, m *User) error

	// Update updates an existing User
	Update(ctx context.Context, m *User) error

	// PartialUpdate updates specific fields of a User
	PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error

	// DeleteByPk deletes a User by primary key
	DeleteByPk(ctx context.Context, pk string) error

	// FindByPk finds a User by primary key
	FindByPk(ctx context.Context, pk string) (*User, error)

	// CreateMany creates multiple User records
	CreateMany(ctx context.Context, models []*User) error

	// UpdateMany updates multiple User records
	UpdateMany(ctx context.Context, models []*User) error

	// DeleteManyByPks deletes multiple User records by primary keys
	DeleteManyByPks(ctx context.Context, pks []string) error

	// FindOne finds a single User with optional where clause and sort expression
	FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*User, error)

	// FindAll finds all User records with optional where clause and sort expression
	FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*User, error)

	// FindPaginated finds User records with pagination, optional where clause and sort expression
	FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*User, error)

	// Count counts User records with optional where clause
	Count(ctx context.Context, where string, args ...interface{}) (int64, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
