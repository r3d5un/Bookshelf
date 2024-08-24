package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type Task struct {
	ID        uuid.UUID  `json:"id"`
	Queue     *string    `json:"queue"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	RunAt     *time.Time `json:"runAt"`
}

type TaskCollection struct {
	CurrentPage  int     `json:"current_page,omitempty"`
	PageSize     int     `json:"page_size,omitempty"`
	FirstPage    int     `json:"first_page,omitempty"`
	LastPage     int     `json:"last_page,omitempty"`
	TotalRecords int     `json:"total_records,omitempty"`
	OrderBy      string  `json:"order_by,omitempty"`
	Data         []*Task `json:"data"`
}

func ReadTask(ctx context.Context, models *data.Models, taskID uuid.UUID) (*Task, error) {
	tq, err := models.TaskQueues.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task := Task{
		ID:        tq.ID,
		Queue:     tq.Queue,
		State:     tq.State,
		CreatedAt: tq.CreatedAt,
		UpdatedAt: tq.UpdatedAt,
		RunAt:     tq.RunAt,
	}

	return &task, nil
}

func ReadAllTasks(
	ctx context.Context,
	models *data.Models,
	filters data.Filters,
) (tc *TaskCollection, err error) {
	tq, metadata, err := models.TaskQueues.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for _, t := range tq {
		task := Task{
			ID:        t.ID,
			Queue:     t.Queue,
			State:     t.State,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			RunAt:     t.RunAt,
		}

		tasks = append(tasks, &task)
	}

	tc = &TaskCollection{
		CurrentPage:  metadata.CurrentPage,
		PageSize:     metadata.PageSize,
		FirstPage:    metadata.FirstPage,
		LastPage:     metadata.LastPage,
		TotalRecords: metadata.TotalRecords,
		Data:         tasks,
	}

	return tc, nil
}
