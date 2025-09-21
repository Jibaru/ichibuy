package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ichibuy/store/internal/domain"
	"strings"
)

type Store = domain.Store

type StoreDAO struct {
	db *sql.DB
}

func NewStoreDAO(db *sql.DB) *StoreDAO {
	return &StoreDAO{db: db}
}

func (dao *StoreDAO) getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("currentTx").(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (dao *StoreDAO) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return dao.db.ExecContext(ctx, query, args...)
}

func (dao *StoreDAO) queryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return dao.db.QueryRowContext(ctx, query, args...)
}

func (dao *StoreDAO) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return dao.db.QueryContext(ctx, query, args...)
}

func (dao *StoreDAO) Create(ctx context.Context, m *Store) error {
	query := `
		INSERT INTO stores (id, name, description, lat, lng, slug, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := dao.execContext(
		ctx,
		query,
		m.ID,
		m.Name,
		m.Description,
		m.Lat,
		m.Lng,
		m.Slug,
		m.UserID,
		m.CreatedAt,
		m.UpdatedAt,
	)

	return err
}

func (dao *StoreDAO) Update(ctx context.Context, m *Store) error {
	query := `
		UPDATE stores
		SET name = $1,
			description = $2,
			lat = $3,
			lng = $4,
			slug = $5,
			user_id = $6,
			created_at = $7,
			updated_at = $8
		WHERE id = $9
	`

	_, err := dao.execContext(ctx, query,
		m.Name,
		m.Description,
		m.Lat,
		m.Lng,
		m.Slug,
		m.UserID,
		m.CreatedAt,
		m.UpdatedAt,
		m.ID,
	)
	return err
}

func (dao *StoreDAO) PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error {
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

	query := fmt.Sprintf(`UPDATE stores SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), i)

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *StoreDAO) DeleteByPk(ctx context.Context, pk string) error {
	query := `DELETE FROM stores WHERE id = $1`
	_, err := dao.execContext(ctx, query, pk)
	return err
}

func (dao *StoreDAO) FindByPk(ctx context.Context, pk string) (*Store, error) {
	query := `
		SELECT id, name, description, lat, lng, slug, user_id, created_at, updated_at
		FROM stores
		WHERE id = $1
	`
	row := dao.queryRowContext(ctx, query, pk)

	var m Store
	err := row.Scan(
		&m.ID,
		&m.Name,
		&m.Description,
		&m.Lat,
		&m.Lng,
		&m.Slug,
		&m.UserID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *StoreDAO) CreateMany(ctx context.Context, models []*Store) error {
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
			model.Lat,
			model.Lng,
			model.Slug,
			model.UserID,
			model.CreatedAt,
			model.UpdatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO stores (id, name, description, lat, lng, slug, user_id, created_at, updated_at)
		VALUES %s
	`, strings.Join(placeholders, ", "))

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *StoreDAO) UpdateMany(ctx context.Context, models []*Store) error {
	if len(models) == 0 {
		return nil
	}

	query := `
		UPDATE stores
		SET name = $1,
			description = $2,
			lat = $3,
			lng = $4,
			slug = $5,
			user_id = $6,
			created_at = $7,
			updated_at = $8
		WHERE id = $9
	`

	for _, model := range models {
		_, err := dao.execContext(ctx, query,
			model.Name,
			model.Description,
			model.Lat,
			model.Lng,
			model.Slug,
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

func (dao *StoreDAO) DeleteManyByPks(ctx context.Context, pks []string) error {
	if len(pks) == 0 {
		return nil
	}

	placeholders := make([]string, len(pks))
	args := make([]interface{}, len(pks))
	for i, pk := range pks {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = pk
	}

	query := fmt.Sprintf(`DELETE FROM stores WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *StoreDAO) FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Store, error) {
	query := `
		SELECT id, name, description, lat, lng, slug, user_id, created_at, updated_at
		FROM stores
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	row := dao.queryRowContext(ctx, query, args...)

	var m Store
	err := row.Scan(
		&m.ID,
		&m.Name,
		&m.Description,
		&m.Lat,
		&m.Lng,
		&m.Slug,
		&m.UserID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *StoreDAO) FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Store, error) {
	query := `
		SELECT id, name, description, lat, lng, slug, user_id, created_at, updated_at
		FROM stores
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

	var models []*Store
	for rows.Next() {
		var m Store
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Lat,
			&m.Lng,
			&m.Slug,
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

func (dao *StoreDAO) FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Store, error) {
	query := `
		SELECT id, name, description, lat, lng, slug, user_id, created_at, updated_at
		FROM stores
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

	var models []*Store
	for rows.Next() {
		var m Store
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Lat,
			&m.Lng,
			&m.Slug,
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

func (dao *StoreDAO) Count(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM stores"

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

func (dao *StoreDAO) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
