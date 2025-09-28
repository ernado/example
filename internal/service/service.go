// Package service implements the business logic of the application.
package service

import (
	"context"

	"github.com/ernado/example"
	"github.com/go-faster/errors"
)

// TODO: Auto-instrument Service.

var _ example.Service = (*Service)(nil)

type Service struct {
	db example.DB
}

func New(db example.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTask(ctx context.Context, title string) (*example.Task, error) {
	task, err := s.db.CreateTask(ctx, title)
	if err != nil {
		return nil, errors.Wrap(err, "create task")
	}
	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id int64) error {
	if err := s.db.DeleteTask(ctx, id); err != nil {
		return errors.Wrap(err, "delete task")
	}
	return nil
}

func (s *Service) ListTasks(ctx context.Context) ([]*example.Task, error) {
	tasks, err := s.db.ListTasks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "list tasks")
	}
	return tasks, nil
}

func (s *Service) GenerateError(ctx context.Context) error {
	if err := s.db.GenerateError(ctx); err != nil {
		return errors.Wrap(err, "generate error")
	}
	return nil
}
