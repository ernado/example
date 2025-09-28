package example

import (
	"context"
	"errors"
)

// Task models a task.
type Task struct {
	ID    int64
	Title string
}

// DB represents a database for tasks.
//
//go:generate go tool moq -fmt goimports -out ./internal/mock/db.go -pkg mock . DB
type DB interface {
	CreateTask(ctx context.Context, title string) (*Task, error)
	ListTasks(ctx context.Context) ([]*Task, error)
	DeleteTask(ctx context.Context, id int64) error
	GenerateError(ctx context.Context) error
}

// Service defines business logic for tasks.
//
//go:generate go tool moq -fmt goimports -out ./internal/mock/service.go -pkg mock . Service
type Service interface {
	CreateTask(ctx context.Context, title string) (*Task, error)
	ListTasks(ctx context.Context) ([]*Task, error)
	DeleteTask(ctx context.Context, id int64) error
	GenerateError(ctx context.Context) error
}

var ErrTaskNotFound = errors.New("task not found")
