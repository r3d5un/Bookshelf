package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskQueueModel(t *testing.T) {
	queue := "test_queue"
	state := data.WaitingTaskState
	timestamp := time.Now()
	tq := data.TaskQueue{
		Queue:     &queue,
		State:     &state,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		RunAt:     &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		insertedTask, err := models.TaskQueues.Insert(
			context.Background(), *tq.Queue, tq.State, tq.RunAt,
		)
		if err != nil {
			t.Errorf("error occurred while inserting new task: %s", err)
			return
		}

		tq = *insertedTask
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.TaskQueues.Get(context.Background(), tq.ID)
		if err != nil {
			t.Errorf("error occurred while reading task: %s", err)
			return
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 100,
			OrderBy:  []string{"id"},
		}
		tasks, metadata, err := models.TaskQueues.GetAll(context.Background(), filters)
		if err != nil {
			t.Errorf("error occurred while reading tasks: %s", err)
			return
		}
		if len(tasks) < 1 {
			t.Errorf("no tasks returned")
			return
		}
		if metadata.CurrentPage != filters.Page {
			t.Errorf("expected page %d in metadata, got %d", filters.Page, metadata.CurrentPage)
			return
		}
	})
}
