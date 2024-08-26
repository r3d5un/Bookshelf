package data_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskQueueModel(t *testing.T) {
	queue := "test_queue"
	state := data.WaitingTaskState
	timestamp := time.Now()
	tq := data.TaskQueue{
		Name:      &queue,
		State:     &state,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		RunAt:     &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		insertedTask, err := models.TaskQueues.Insert(context.Background(), tq)
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

	t.Run("Update", func(t *testing.T) {
		newTaskState := data.RunningTaskState
		tq.State = &newTaskState

		updatedTask, err := models.TaskQueues.Update(context.Background(), tq)
		if err != nil {
			t.Errorf("error occurred while updating task: %s", err)
			return
		}
		if *updatedTask.State != newTaskState {
			t.Errorf("expected task state %s, got %s", newTaskState, *updatedTask.State)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := models.TaskQueues.Delete(context.Background(), tq.ID)
		if err != nil {
			t.Errorf("error occurred while deleting task: %s\n", err)
			return
		}
	})

	t.Run("ConsumeByID", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		taskName := "test_queue"
		timestamp := time.Now().Add(-1 * time.Hour)
		state := data.WaitingTaskState
		task := data.TaskQueue{
			Name:     &taskName,
			State:    &state,
			RunAt:    &timestamp,
			TaskData: nil,
		}
		targetTask, err := models.TaskQueues.Insert(ctx, task)
		if err != nil {
			t.Fatalf("error inserting task: %s", err)
		}

		taskCh := make(chan data.TaskQueue, 1)
		taskRunResultCh := make(chan error, 1)

		go func() {
			// Ensure this happens after ConsumeByID has sent the task
			task := <-taskCh
			slog.Info("task retrieved", "task", task)

			// Simulate task processing completion
			taskRunResultCh <- nil
			close(taskRunResultCh)
		}()

		err = models.TaskQueues.ConsumeByID(
			ctx,
			taskCh,
			taskRunResultCh,
			targetTask.ID,
		)
		if err != nil {
			t.Errorf("error occurred while consuming task by ID: %s\n", err)
			return
		}

		close(taskCh)
	})
}
