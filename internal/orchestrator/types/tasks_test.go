package types_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

func TestTaskTypes(t *testing.T) {
	_ = []types.Task{
		types.NewTask("task1", "* * * * *", false, time.Now()),
		types.NewTask("task2", "* * * * *", false, time.Now()),
		types.NewTask("task3", "* * * * *", false, time.Now()),
		types.NewTask("task4", "* * * * *", false, time.Now()),
	}

	t.Run("CreateTask", func(t *testing.T) {
		_, err := types.CreateTask(
			context.Background(),
			models,
			types.NewTask("task0", "* * * * *", false, time.Now()),
		)
		if err != nil {
			t.Errorf("an error occurred while creating a new task: %s\n", err)
			return
		}
	})
}
