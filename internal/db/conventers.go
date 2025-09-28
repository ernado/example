package db

import (
	"github.com/ernado/example"
	"github.com/ernado/example/internal/ent"
)

func taskFromDB(task *ent.Task) *example.Task {
	return &example.Task{
		Title: task.Title,
		ID:    int64(task.ID),
	}
}

func tasksFromDB(tasks []*ent.Task) []*example.Task {
	result := make([]*example.Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, taskFromDB(task))
	}
	return result
}
