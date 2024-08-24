package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskQueueModel(t *testing.T) {
	queue := "test_queue"
	state := data.WaitingTaskState
	timestamp := time.Now()
	tq := data.TaskQueue{
		ID:        uuid.New(),
		Queue:     &queue,
		State:     &state,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		RunAt:     &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.TaskQueues.Insert(
			context.Background(), *tq.Queue, tq.State, tq.RunAt,
		)
		if err != nil {
			t.Errorf("error occurred while inserting new task: %s", err)
			return
		}
	})
}
