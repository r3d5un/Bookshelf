package types_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

func TestTaskTypes(t *testing.T) {
	tasks := []types.Task{
		types.NewTask("task1", "* * * * *", false, time.Now()),
		types.NewTask("task2", "* * * * *", false, time.Now()),
		types.NewTask("task3", "* * * * *", false, time.Now()),
		types.NewTask("task4", "* * * * *", false, time.Now()),
	}

	filters := data.Filters{
		Page:     1,
		PageSize: 50_000,
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

	t.Run("ReadAllTasks", func(t *testing.T) {
		taskCollection, err := types.ReadAllTasks(context.Background(), models, filters)
		if err != nil {
			t.Errorf("an error occurred while reading all tasks: %s\n", err)
			return
		}
		if len(taskCollection.Tasks) < 1 {
			t.Errorf("no tasks returned")
			return
		}
	})

	t.Run("UpdateTask", func(t *testing.T) {
		cronExpr := "1 * * * *"
		task.CronExpr = &cronExpr
		updatedTask, err := types.UpdateTask(context.Background(), models, task)
		if err != nil {
			t.Errorf("an error occurred while updating task: %s\n", err)
			return
		}
		if *updatedTask.CronExpr != cronExpr {
			t.Errorf("expected %s, got %s\n", cronExpr, *updatedTask.CronExpr)
			return
		}
	})

	t.Run("SyncTasks", func(t *testing.T) {
		err := types.SyncTasks(context.Background(), models, tasks)
		if err != nil {
			t.Errorf("an error occurred while syncing tasks: %s\n", err)
			return
		}

		syncedTasks, err := types.ReadAllTasks(context.Background(), models, filters)
		if err != nil {
			t.Errorf("an error occurred while reading synced tasks: %s\n", err)
			return
		}

		if len(syncedTasks.Tasks) < len(tasks) {
			t.Errorf(
				"expected equal or more than %d tasks, got %d\n",
				len(tasks), len(syncedTasks.Tasks),
			)
			return
		}

		_, err = types.ReadTask(context.Background(), models, task.Name)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				return
			default:
				t.Errorf("unable to read desynced task: %s\n", err)
			}
			return
		}
	})
}
