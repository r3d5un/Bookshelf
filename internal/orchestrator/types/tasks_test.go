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

	var task types.Task

	t.Run("CreateTask", func(t *testing.T) {
		createdTask, err := types.CreateTask(
			context.Background(),
			models,
			types.NewTask("task0", "* * * * *", false, time.Now()),
		)
		if err != nil {
			t.Errorf("an error occurred while creating a new task: %s\n", err)
			return
		}

		task = *createdTask
	})

	t.Run("ReadTask", func(t *testing.T) {
		readTask, err := types.ReadTask(context.Background(), models, task.Name)
		if err != nil {
			t.Errorf("an error occurred while readin task: %s\n", err)
			return
		}
		if readTask.Name != task.Name {
			t.Errorf("expected task name %s, got %s\n", task.Name, readTask.Name)
			return
		}
	})
}
