package handler

import (
	"github.com/ernado/example"
	"github.com/ernado/example/internal/oas"
)

func convertToOASTask(t *example.Task) *oas.Task {
	return &oas.Task{
		ID:    t.ID,
		Title: t.Title,
	}
}

func convertToOASTasks(ts []*example.Task) []oas.Task {
	res := make([]oas.Task, 0, len(ts))
	for _, t := range ts {
		res = append(res, *convertToOASTask(t))
	}
	return res
}
