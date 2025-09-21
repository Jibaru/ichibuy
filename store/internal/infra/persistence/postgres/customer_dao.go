package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ichibuy/store/internal/domain"
	"strings"
)

type Customer = domain.Customer

type CustomerDAO struct {
	db *sql.DB
}

func NewCustomerDAO(db *sql.DB) *CustomerDAO {
	return &CustomerDAO{db: db}
}

func (dao *CustomerDAO) getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("currentTx").(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (dao *CustomerDAO) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return dao.db.ExecContext(ctx, query, args...)
}

func (dao *CustomerDAO) queryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return dao.db.QueryRowContext(ctx, query, args...)
}

func (dao *CustomerDAO) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return dao.db.QueryContext(ctx, query, args...)
}

func (dao *CustomerDAO) Create(ctx context.Context, m *Customer) error {
	query := `
		INSERT INTO customers (id, first_name, last_name, email, phone, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := dao.execContext(
		ctx,
		query,
		m.ID,
		m.FirstName,
		m.LastName,
		m.Email,
		m.Phone,
		m.UserID,
		m.CreatedAt,
		m.UpdatedAt,
	)

	return err
}

func (dao *CustomerDAO) Update(ctx context.Context, m *Customer) error {
	query := `
		UPDATE customers
		SET first_name = $1,
			last_name = $2,
			email = $3,
			phone = $4,
			user_id = $5,
			created_at = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := dao.execContext(ctx, query,
		m.FirstName,
		m.LastName,
		m.Email,
		m.Phone,
		m.UserID,
		m.CreatedAt,
		m.UpdatedAt,
		m.ID,
	)
	return err
}

func (dao *CustomerDAO) PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error {
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

	query := fmt.Sprintf(`UPDATE customers SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), i)

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *CustomerDAO) DeleteByPk(ctx context.Context, pk string) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := dao.execContext(ctx, query, pk)
	return err
}

func (dao *CustomerDAO) FindByPk(ctx context.Context, pk string) (*Customer, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, user_id, created_at, updated_at
		FROM customers
		WHERE id = $1
	`
	row := dao.queryRowContext(ctx, query, pk)

	var m Customer
	err := row.Scan(
		&m.ID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Phone,
		&m.UserID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *CustomerDAO) CreateMany(ctx context.Context, models []*Customer) error {
	if len(models) == 0 {
		return nil
	}

	placeholders := make([]string, len(models))
	args := make([]interface{}, 0, len(models)*8)

	for i, model := range models {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)

		args = append(args,
			model.ID,
			model.FirstName,
			model.LastName,
			model.Email,
			model.Phone,
			model.UserID,
			model.CreatedAt,
			model.UpdatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO customers (id, first_name, last_name, email, phone, user_id, created_at, updated_at)
		VALUES %s
	`, strings.Join(placeholders, ", "))

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *CustomerDAO) UpdateMany(ctx context.Context, models []*Customer) error {
	if len(models) == 0 {
		return nil
	}

	query := `
		UPDATE customers
		SET first_name = $1,
			last_name = $2,
			email = $3,
			phone = $4,
			user_id = $5,
			created_at = $6,
			updated_at = $7
		WHERE id = $8
	`

	for _, model := range models {
		_, err := dao.execContext(ctx, query,
			model.FirstName,
			model.LastName,
			model.Email,
			model.Phone,
			model.UserID,
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

func (dao *CustomerDAO) DeleteManyByPks(ctx context.Context, pks []string) error {
	if len(pks) == 0 {
		return nil
	}

	placeholders := make([]string, len(pks))
	args := make([]interface{}, len(pks))
	for i, pk := range pks {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = pk
	}

	query := fmt.Sprintf(`DELETE FROM customers WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *CustomerDAO) FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Customer, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, user_id, created_at, updated_at
		FROM customers
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	row := dao.queryRowContext(ctx, query, args...)

	var m Customer
	err := row.Scan(
		&m.ID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Phone,
		&m.UserID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *CustomerDAO) FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Customer, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, user_id, created_at, updated_at
		FROM customers
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

	var models []*Customer
	for rows.Next() {
		var m Customer
		err := rows.Scan(
			&m.ID,
			&m.FirstName,
			&m.LastName,
			&m.Email,
			&m.Phone,
			&m.UserID,
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

func (dao *CustomerDAO) FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Customer, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, user_id, created_at, updated_at
		FROM customers
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

	var models []*Customer
	for rows.Next() {
		var m Customer
		err := rows.Scan(
			&m.ID,
			&m.FirstName,
			&m.LastName,
			&m.Email,
			&m.Phone,
			&m.UserID,
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

func (dao *CustomerDAO) Count(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM customers"

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

func (dao *CustomerDAO) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
