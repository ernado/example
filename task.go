package example

import (
	"context"
	"errors"
)

type Task struct {
	ID    int64
	Title string
}

type DB interface {
	CreateTask(ctx context.Context, task *Task) error
	ListTasks(ctx context.Context) ([]*Task, error)
	DeleteTask(ctx context.Context, id int64) error
}

var ErrTaskNotFound = errors.New("task not found")
