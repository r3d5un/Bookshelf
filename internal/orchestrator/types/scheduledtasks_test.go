package types_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

func TestScheduledTaskTypes(t *testing.T) {
	insertedTask, err := types.CreateTask(
		context.Background(),
		models,
		types.NewTask("test_queue", "* * * * *", false, time.Now(), nil),
	)
	if err != nil {
		t.Errorf("an error occurred while creating parent task for overview: %s\n", err)
		return
	}

	state := string(data.WaitingTaskState)
	timestamp := time.Now()
	tq := types.ScheduledTask{
		Name:      &insertedTask.Name,
		State:     &state,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		RunAt:     &timestamp,
	}

	t.Run("CreateScheduledTask", func(t *testing.T) {
		createdTask, err := types.ScheduleTask(context.Background(), models, tq)
		if err != nil {
			t.Errorf("error occurred while creating task: %s\n", err)
			return
		}

		tq = *createdTask
	})

	t.Run("ReadScheduledTask", func(t *testing.T) {
		_, err := types.ReadScheduledTask(context.Background(), models, tq.ID)
		if err != nil {
			t.Errorf("error occurred while reading task: %s\n", err)
		}
	})

	t.Run("ReadAllScheduledTask", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 100,
			OrderBy:  []string{"id"},
		}
		tc, err := types.ReadAllScheudledTasks(context.Background(), models, filters)
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

	t.Run("UpdateScheduledTask", func(t *testing.T) {
		newTaskState := string(data.RunningTaskState)
		tq.State = &newTaskState

		updatedTask, err := types.UpdateScheduledTask(context.Background(), models, tq)
		if err != nil {
			t.Errorf("error occurred while creating task: %s\n", err)
			return
		}
		if *updatedTask.State != newTaskState {
			t.Errorf("expected task state %s, got %s", newTaskState, *updatedTask.State)
		}
	})

	t.Run("DeleteScheduledTask", func(t *testing.T) {
		_, err := types.DeleteScheduledTask(context.Background(), models, tq.ID)
		if err != nil {
			t.Errorf("error occurred while creating task: %s\n", err)
			return
		}
	})
}
