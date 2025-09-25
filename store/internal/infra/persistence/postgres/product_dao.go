package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ichibuy/store/internal/domain"
	"strings"
)

type Product = domain.Product

type ProductDAO struct {
	db *sql.DB
}

func NewProductDAO(db *sql.DB) *ProductDAO {
	return &ProductDAO{db: db}
}

func (dao *ProductDAO) getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("currentTx").(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (dao *ProductDAO) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return dao.db.ExecContext(ctx, query, args...)
}

func (dao *ProductDAO) queryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return dao.db.QueryRowContext(ctx, query, args...)
}

func (dao *ProductDAO) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return dao.db.QueryContext(ctx, query, args...)
}

func (dao *ProductDAO) Create(ctx context.Context, m *Product) error {
	query := `
		INSERT INTO products (id, name, description, active, store_id, images, prices, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := dao.execContext(
		ctx,
		query,
		m.ID,
		m.Name,
		m.Description,
		m.Active,
		m.StoreID,
		m.Images,
		m.Prices,
		m.CreatedAt,
		m.UpdatedAt,
	)

	return err
}

func (dao *ProductDAO) Update(ctx context.Context, m *Product) error {
	query := `
		UPDATE products
		SET name = $1,
			description = $2,
			active = $3,
			store_id = $4,
			images = $5,
			prices = $6,
			created_at = $7,
			updated_at = $8
		WHERE id = $9
	`

	_, err := dao.execContext(ctx, query,
		m.Name,
		m.Description,
		m.Active,
		m.StoreID,
		m.Images,
		m.Prices,
		m.CreatedAt,
		m.UpdatedAt,
		m.ID,
	)
	return err
}

func (dao *ProductDAO) PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields)+1)
	i := 1

	for field, value := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, i))
		args = append(args, value)
		i++
	}

	args = append(args, pk)

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), i)

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *ProductDAO) DeleteByPk(ctx context.Context, pk string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := dao.execContext(ctx, query, pk)
	return err
}

func (dao *ProductDAO) FindByPk(ctx context.Context, pk string) (*Product, error) {
	query := `
		SELECT id, name, description, active, store_id, images, prices, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	row := dao.queryRowContext(ctx, query, pk)

	var m Product
	err := row.Scan(
		&m.ID,
		&m.Name,
		&m.Description,
		&m.Active,
		&m.StoreID,
		&m.Images,
		&m.Prices,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *ProductDAO) CreateMany(ctx context.Context, models []*Product) error {
	if len(models) == 0 {
		return nil
	}

	placeholders := make([]string, len(models))
	args := make([]interface{}, 0, len(models)*9)

	for i, model := range models {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9)

		args = append(args,
			model.ID,
			model.Name,
			model.Description,
			model.Active,
			model.StoreID,
			model.Images,
			model.Prices,
			model.CreatedAt,
			model.UpdatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO products (id, name, description, active, store_id, images, prices, created_at, updated_at)
		VALUES %s
	`, strings.Join(placeholders, ", "))

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *ProductDAO) UpdateMany(ctx context.Context, models []*Product) error {
	if len(models) == 0 {
		return nil
	}

	query := `
		UPDATE products
		SET name = $1,
			description = $2,
			active = $3,
			store_id = $4,
			images = $5,
			prices = $6,
			created_at = $7,
			updated_at = $8
		WHERE id = $9
	`

	for _, model := range models {
		_, err := dao.execContext(ctx, query,
			model.Name,
			model.Description,
			model.Active,
			model.StoreID,
			model.Images,
			model.Prices,
			model.CreatedAt,
			model.UpdatedAt,
			model.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dao *ProductDAO) DeleteManyByPks(ctx context.Context, pks []string) error {
	if len(pks) == 0 {
		return nil
	}

	placeholders := make([]string, len(pks))
	args := make([]interface{}, len(pks))
	for i, pk := range pks {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = pk
	}

	query := fmt.Sprintf(`DELETE FROM products WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *ProductDAO) FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Product, error) {
	query := `
		SELECT id, name, description, active, store_id, images, prices, created_at, updated_at
		FROM products
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	row := dao.queryRowContext(ctx, query, args...)

	var m Product
	err := row.Scan(
		&m.ID,
		&m.Name,
		&m.Description,
		&m.Active,
		&m.StoreID,
		&m.Images,
		&m.Prices,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *ProductDAO) FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Product, error) {
	query := `
		SELECT id, name, description, active, store_id, images, prices, created_at, updated_at
		FROM products
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	rows, err := dao.queryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*Product
	for rows.Next() {
		var m Product
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Active,
			&m.StoreID,
			&m.Images,
			&m.Prices,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (dao *ProductDAO) FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Product, error) {
	query := `
		SELECT id, name, description, active, store_id, images, prices, created_at, updated_at
		FROM products
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := dao.queryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*Product
	for rows.Next() {
		var m Product
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Active,
			&m.StoreID,
			&m.Images,
			&m.Prices,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (dao *ProductDAO) Count(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM products"

	if where != "" {
		query += " WHERE " + where
	}

	row := dao.queryRowContext(ctx, query, args...)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *ProductDAO) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := dao.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, "currentTx", tx)

	err = fn(ctxWithTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
