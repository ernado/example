package entdb

import (
	"context"

	"github.com/ernado/example"
	"github.com/ernado/example/internal/ent"
	"github.com/pkg/errors"
)

var _ example.DB = (*DB)(nil)

type DB struct {
	ent *ent.Client
}

func (db DB) CreateTask(ctx context.Context, task *example.Task) error {
	_, err := db.ent.Task.Create().
		SetTitle(task.Title).
		Save(ctx)
	if err != nil {
		return errors.Wrap(err, "create task")
	}

	return nil
}

func (db DB) ListTasks(ctx context.Context) ([]*example.Task, error) {
	tasksDB, err := db.ent.Task.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "query tasks")
	}
	tasks := tasksFromDB(tasksDB)

	return tasks, nil
}

func (db DB) DeleteTask(ctx context.Context, id int64) error {
	err := db.ent.Task.DeleteOneID(int(id)).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return example.ErrTaskNotFound
		}
		return errors.Wrap(err, "delete task")
	}

	return nil
}

func New(ent *ent.Client) *DB {
	return &DB{
		ent: ent,
	}
}
