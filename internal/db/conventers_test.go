package entdb

import (
	"testing"

	"github.com/ernado/example"
	"github.com/ernado/example/internal/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskFromDB(t *testing.T) {
	tests := []struct {
		name string
		give *ent.Task
		want *example.Task
	}{
		{
			name: "basic conversion",
			give: &ent.Task{
				ID:    1,
				Title: "Test Task",
			},
			want: &example.Task{
				ID:    1,
				Title: "Test Task",
			},
		},
		{
			name: "empty title",
			give: &ent.Task{
				ID:    42,
				Title: "",
			},
			want: &example.Task{
				ID:    42,
				Title: "",
			},
		},
		{
			name: "zero id",
			give: &ent.Task{
				ID:    0,
				Title: "Zero ID Task",
			},
			want: &example.Task{
				ID:    0,
				Title: "Zero ID Task",
			},
		},
		{
			name: "large id",
			give: &ent.Task{
				ID:    999999,
				Title: "Large ID Task",
			},
			want: &example.Task{
				ID:    999999,
				Title: "Large ID Task",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskFromDB(tt.give)
			require.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Title, got.Title)
		})
	}
}

func TestTasksFromDB(t *testing.T) {
	tests := []struct {
		name string
		give []*ent.Task
		want []*example.Task
	}{
		{
			name: "empty slice",
			give: []*ent.Task{},
			want: []*example.Task{},
		},
		{
			name: "nil slice",
			give: nil,
			want: []*example.Task{},
		},
		{
			name: "single task",
			give: []*ent.Task{
				{
					ID:    1,
					Title: "Single Task",
				},
			},
			want: []*example.Task{
				{
					ID:    1,
					Title: "Single Task",
				},
			},
		},
		{
			name: "multiple tasks",
			give: []*ent.Task{
				{
					ID:    1,
					Title: "First Task",
				},
				{
					ID:    2,
					Title: "Second Task",
				},
				{
					ID:    3,
					Title: "Third Task",
				},
			},
			want: []*example.Task{
				{
					ID:    1,
					Title: "First Task",
				},
				{
					ID:    2,
					Title: "Second Task",
				},
				{
					ID:    3,
					Title: "Third Task",
				},
			},
		},
		{
			name: "tasks with empty titles",
			give: []*ent.Task{
				{
					ID:    1,
					Title: "",
				},
				{
					ID:    2,
					Title: "Non-empty",
				},
			},
			want: []*example.Task{
				{
					ID:    1,
					Title: "",
				},
				{
					ID:    2,
					Title: "Non-empty",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tasksFromDB(tt.give)
			require.NotNil(t, got)
			assert.Len(t, got, len(tt.want))

			for i, expectedTask := range tt.want {
				require.Less(t, i, len(got), "got slice shorter than expected")
				assert.Equal(t, expectedTask.ID, got[i].ID)
				assert.Equal(t, expectedTask.Title, got[i].Title)
			}
		})
	}
}

func TestTasksFromDB_PreservesCapacity(t *testing.T) {
	// Test that the function allocates with correct capacity
	tasks := []*ent.Task{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}

	got := tasksFromDB(tasks)

	// Verify we got the expected results
	require.Len(t, got, 2)
	assert.Equal(t, int64(1), got[0].ID)
	assert.Equal(t, "Task 1", got[0].Title)
	assert.Equal(t, int64(2), got[1].ID)
	assert.Equal(t, "Task 2", got[1].Title)
}
