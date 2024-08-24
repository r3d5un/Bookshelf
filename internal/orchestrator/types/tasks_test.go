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
}
