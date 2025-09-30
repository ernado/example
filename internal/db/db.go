package db

import (
	"context"

	"github.com/ernado/example"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ example.DB = (*DB)(nil)

type DB struct {
	pgx *pgxpool.Pool
}

func (db DB) CreateTask(ctx context.Context, title string) (*example.Task, error) {
	const query = `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id
	`
	var id int64
	err := db.pgx.QueryRow(ctx, query, title).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &example.Task{
		ID:    id,
		Title: title,
	}, nil
}

func (db DB) ListTasks(ctx context.Context) ([]*example.Task, error) {
	const query = `
		SELECT id, title
		FROM tasks
		ORDER BY id
	`

	rows, err := db.pgx.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query tasks")
	}
	defer rows.Close()

	var tasks []*example.Task
	for rows.Next() {
		var t example.Task
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil, errors.Wrap(err, "scan row")
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterate rows")
	}

	return tasks, nil
}

func (db DB) DeleteTask(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM tasks
		WHERE id = $1
	`
	cmdTag, err := db.pgx.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "exec delete")
	}
	if cmdTag.RowsAffected() == 0 {
		return example.ErrTaskNotFound
	}

	return nil
}

func (db DB) GenerateError(ctx context.Context) error {
	return errors.New("generated error")
}

func New(pgx *pgxpool.Pool) *DB {
	return &DB{pgx: pgx}
}
