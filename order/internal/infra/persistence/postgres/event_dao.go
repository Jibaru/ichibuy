package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ichibuy/order/internal/domain"
	"strings"
)

type Event = domain.Event

type EventDAO struct {
	db *sql.DB
}

func NewEventDAO(db *sql.DB) *EventDAO {
	return &EventDAO{db: db}
}

func (dao *EventDAO) getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("currentTx").(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (dao *EventDAO) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return dao.db.ExecContext(ctx, query, args...)
}

func (dao *EventDAO) queryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return dao.db.QueryRowContext(ctx, query, args...)
}

func (dao *EventDAO) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := dao.getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return dao.db.QueryContext(ctx, query, args...)
}

func (dao *EventDAO) Create(ctx context.Context, m *Event) error {
	query := `
		INSERT INTO events (id, type, data, timestamp)
		VALUES ($1, $2, $3, $4)
	`

	_, err := dao.execContext(
		ctx,
		query,
		m.ID,
		m.Type,
		m.Data,
		m.Timestamp,
	)

	return err
}

func (dao *EventDAO) Update(ctx context.Context, m *Event) error {
	query := `
		UPDATE events
		SET type = $1,
			data = $2,
			timestamp = $3
		WHERE id = $4
	`

	_, err := dao.execContext(ctx, query,
		m.Type,
		m.Data,
		m.Timestamp,
		m.ID,
	)
	return err
}

func (dao *EventDAO) PartialUpdate(ctx context.Context, pk string, fields map[string]interface{}) error {
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

	query := fmt.Sprintf(`UPDATE events SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), i)

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *EventDAO) DeleteByPk(ctx context.Context, pk string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := dao.execContext(ctx, query, pk)
	return err
}

func (dao *EventDAO) FindByPk(ctx context.Context, pk string) (*Event, error) {
	query := `
		SELECT id, type, data, timestamp
		FROM events
		WHERE id = $1
	`
	row := dao.queryRowContext(ctx, query, pk)

	var m Event
	err := row.Scan(
		&m.ID,
		&m.Type,
		&m.Data,
		&m.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *EventDAO) CreateMany(ctx context.Context, models []*Event) error {
	if len(models) == 0 {
		return nil
	}

	placeholders := make([]string, len(models))
	args := make([]interface{}, 0, len(models)*4)

	for i, model := range models {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4)

		args = append(args,
			model.ID,
			model.Type,
			model.Data,
			model.Timestamp,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO events (id, type, data, timestamp)
		VALUES %s
	`, strings.Join(placeholders, ", "))

	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *EventDAO) UpdateMany(ctx context.Context, models []*Event) error {
	if len(models) == 0 {
		return nil
	}

	query := `
		UPDATE events
		SET type = $1,
			data = $2,
			timestamp = $3
		WHERE id = $4
	`

	for _, model := range models {
		_, err := dao.execContext(ctx, query,
			model.Type,
			model.Data,
			model.Timestamp,
			model.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dao *EventDAO) DeleteManyByPks(ctx context.Context, pks []string) error {
	if len(pks) == 0 {
		return nil
	}

	placeholders := make([]string, len(pks))
	args := make([]interface{}, len(pks))
	for i, pk := range pks {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = pk
	}

	query := fmt.Sprintf(`DELETE FROM events WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err := dao.execContext(ctx, query, args...)
	return err
}

func (dao *EventDAO) FindOne(ctx context.Context, where string, sort string, args ...interface{}) (*Event, error) {
	query := `
		SELECT id, type, data, timestamp
		FROM events
	`

	if where != "" {
		query += " WHERE " + where
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	row := dao.queryRowContext(ctx, query, args...)

	var m Event
	err := row.Scan(
		&m.ID,
		&m.Type,
		&m.Data,
		&m.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *EventDAO) FindAll(ctx context.Context, where string, sort string, args ...interface{}) ([]*Event, error) {
	query := `
		SELECT id, type, data, timestamp
		FROM events
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

	var models []*Event
	for rows.Next() {
		var m Event
		err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Data,
			&m.Timestamp,
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

func (dao *EventDAO) FindPaginated(ctx context.Context, limit, offset int, where string, sort string, args ...interface{}) ([]*Event, error) {
	query := `
		SELECT id, type, data, timestamp
		FROM events
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

	var models []*Event
	for rows.Next() {
		var m Event
		err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Data,
			&m.Timestamp,
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

func (dao *EventDAO) Count(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM events"

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

func (dao *EventDAO) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
