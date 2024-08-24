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
