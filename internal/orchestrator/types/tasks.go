package types

import (
	"context"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type Task struct {
	Name      string     `json:"name"`
	CronExpr  *string    `json:"cronExpr,omitempty"`
	Enabled   *bool      `json:"enabled,omitempty"`
	Deleted   *bool      `json:"deleted,omitempty"`
	UpdatedAt *time.Time `json:"timestamp,omitempty"`
}

func ReadTask(ctx context.Context, models *data.Models, taskName string) (*Task, error) {
	taskRow, err := models.Tasks.Get(ctx, taskName)
	if err != nil {
		return nil, err
	}

	task := Task{
		Name:      taskRow.Name.String,
		CronExpr:  &taskRow.CronExpr.String,
		Enabled:   &taskRow.Enabled.Bool,
		Deleted:   &taskRow.Deleted.Bool,
		UpdatedAt: &taskRow.UpdatedAt.Time,
	}

	return &task, nil
}
