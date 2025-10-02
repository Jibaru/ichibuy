package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ichibuy/order/internal/domain"
	"strings"
)

type Order = domain.Order

type OrderDAO struct {
	db *sql.DB
}

func NewOrderDAO(db *sql.DB) *OrderDAO {
	return &OrderDAO{db: db}
}

func (dao *OrderDAO) getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("currentTx").(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (dao *OrderDAO) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return dao.db.ExecContext(ctx, query, args...)
}

func (dao *OrderDAO) queryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return dao.db.QueryRowContext(ctx, query, args...)
}

func (dao *OrderDAO) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return dao.db.QueryContext(ctx, query, args...)
}

func (dao *OrderDAO) Create(ctx context.Context, m *Order) error {
	query := `
		INSERT INTO orders (id, code, current_status, order_lines, customer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := dao.execContext(
		ctx,
		query,
		m.ID,
		m.Code,
		m.CurrentStatus,
		m.OrderLines,
		m.CustomerID,
		m.CreatedAt,
		m.UpdatedAt,
	)

	return err
}

func (dao *OrderDAO) Update(ctx context.Context, m *Order) error {
	query := `
		UPDATE orders
		SET code = $1,
			current_status = $2,
			order_lines = $3,
			customer_id = $4,
			created_at = $5,
			updated_at = $6
		WHERE id = $7
	`

	_, err := dao.execContext(ctx, query,
		m.Code,
		m.CurrentStatus,
		m.OrderLines,
		m.CustomerID,
		m.CreatedAt,
		m.UpdatedAt,
		m.ID,
	)
	return err
}

func (dao *OrderDAO) PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error {
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

	query := fmt.Sprintf(`UPDATE orders SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), i)

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *OrderDAO) DeleteByPk(ctx context.Context, pk string) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := dao.execContext(ctx, query, pk)
	return err
}

func (dao *OrderDAO) FindByPk(ctx context.Context, pk string) (*Order, error) {
	query := `
		SELECT id, code, current_status, order_lines, customer_id, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	row := dao.queryRowContext(ctx, query, pk)

	var m Order
	err := row.Scan(
		&m.ID,
		&m.Code,
		&m.CurrentStatus,
		&m.OrderLines,
		&m.CustomerID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *OrderDAO) CreateMany(ctx context.Context, models []*Order) error {
	if len(models) == 0 {
		return nil
	}

	placeholders := make([]string, len(models))
	args := make([]interface{}, 0, len(models)*7)

	for i, model := range models {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7)

		args = append(args,
			model.ID,
			model.Code,
			model.CurrentStatus,
			model.OrderLines,
			model.CustomerID,
			model.CreatedAt,
			model.UpdatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO orders (id, code, current_status, order_lines, customer_id, created_at, updated_at)
		VALUES %s
	`, strings.Join(placeholders, ", "))

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *OrderDAO) UpdateMany(ctx context.Context, models []*Order) error {
	if len(models) == 0 {
		return nil
	}

	query := `
		UPDATE orders
		SET code = $1,
			current_status = $2,
			order_lines = $3,
			customer_id = $4,
			created_at = $5,
			updated_at = $6
		WHERE id = $7
	`

	for _, model := range models {
		_, err := dao.execContext(ctx, query,
			model.Code,
			model.CurrentStatus,
			model.OrderLines,
			model.CustomerID,
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

func (dao *OrderDAO) DeleteManyByPks(ctx context.Context, pks []string) error {
	if len(pks) == 0 {
		return nil
	}

	placeholders := make([]string, len(pks))
	args := make([]interface{}, len(pks))
	for i, pk := range pks {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = pk
	}

	query := fmt.Sprintf(`DELETE FROM orders WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *OrderDAO) FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Order, error) {
	query := `
		SELECT id, code, current_status, order_lines, customer_id, created_at, updated_at
		FROM orders
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	row := dao.queryRowContext(ctx, query, args...)

	var m Order
	err := row.Scan(
		&m.ID,
		&m.Code,
		&m.CurrentStatus,
		&m.OrderLines,
		&m.CustomerID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *OrderDAO) FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Order, error) {
	query := `
		SELECT id, code, current_status, order_lines, customer_id, created_at, updated_at
		FROM orders
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

	var models []*Order
	for rows.Next() {
		var m Order
		err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.CurrentStatus,
			&m.OrderLines,
			&m.CustomerID,
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

func (dao *OrderDAO) FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Order, error) {
	query := `
		SELECT id, code, current_status, order_lines, customer_id, created_at, updated_at
		FROM orders
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

	var models []*Order
	for rows.Next() {
		var m Order
		err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.CurrentStatus,
			&m.OrderLines,
			&m.CustomerID,
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

func (dao *OrderDAO) Count(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM orders"

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

func (dao *OrderDAO) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
