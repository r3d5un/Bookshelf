package types_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

func TestTaskTypes(t *testing.T) {
	queue := "test_queue"
	state := data.WaitingTaskState
	timestamp := time.Now()
	tq := types.Task{
		Queue:     &queue,
		State:     &state,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		RunAt:     &timestamp,
	}

	t.Run("CreateTask", func(t *testing.T) {
		createdTask, err := types.CreateTask(context.Background(), models, tq)
		if err != nil {
			t.Errorf("error occurred while creating task: %s\n", err)
			return
		}

		tq = *createdTask
	})

	t.Run("ReadTask", func(t *testing.T) {
		_, err := types.ReadTask(context.Background(), models, tq.ID)
		if err != nil {
			t.Errorf("error occurred while reading task: %s\n", err)
		}
	})

	t.Run("ReadAllTasks", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 100,
			OrderBy:  []string{"id"},
		}
		tc, err := types.ReadAllTasks(context.Background(), models, filters)
		if err != nil {
			t.Errorf("error occurred while reading tasks: %s\n", err)
		}
		if len(tc.Data) < 1 {
			t.Errorf("no tasks returned")
			return
		}
		if tc.CurrentPage != filters.Page {
			t.Errorf("expected page %d in metadata, got %d", filters.Page, tc.CurrentPage)
			return
		}
	})
}
