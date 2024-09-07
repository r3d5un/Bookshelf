package data_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskQueueModel(t *testing.T) {
	task := data.Task{
		Name:      sql.NullString{String: "test_queue", Valid: true},
		CronExpr:  sql.NullString{String: "* * * * *", Valid: true},
		Enabled:   sql.NullBool{Bool: false, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	_, err := models.Tasks.Insert(context.Background(), task)
	if err != nil {
		t.Errorf("error occurred while inserting task: %s\n", err)
		return
	}

	state := string(data.WaitingTaskState)
	timestamp := time.Now()
	tq := data.TaskQueue{
		Name:      &task.Name.String,
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
		newTaskState := string(data.RunningTaskState)
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
}
