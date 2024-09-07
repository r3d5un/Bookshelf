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
	UpdatedAt *time.Time `json:"timestamp,omitempty"`
}

type TaskCollection struct {
	CurrentPage  int     `json:"current_page,omitempty"`
	PageSize     int     `json:"page_size,omitempty"`
	FirstPage    int     `json:"first_page,omitempty"`
	LastPage     int     `json:"last_page,omitempty"`
	TotalRecords int     `json:"total_records,omitempty"`
	OrderBy      string  `json:"order_by,omitempty"`
	Tasks        []*Task `json:"tasks"`
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
		UpdatedAt: &taskRow.UpdatedAt.Time,
	}

	return &task, nil
}

func ReadAllTasks(
	ctx context.Context,
	models *data.Models,
	filters data.Filters,
) (*TaskCollection, error) {
	// CreateOrderBy clause uses id as default value, which will cause an error for the TaskModel
	// since it doesn't have a id field. It uses the name column as the primary key.
	if len(filters.OrderBy) < 1 {
		filters.OrderBy = []string{"name"}
	}

	taskRows, metadata, err := models.Tasks.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for _, t := range taskRows {
		task := Task{
			Name:      t.Name.String,
			CronExpr:  &t.CronExpr.String,
			Enabled:   &t.Enabled.Bool,
			UpdatedAt: &t.UpdatedAt.Time,
		}

		tasks = append(tasks, &task)
	}

	tc := &TaskCollection{
		CurrentPage:  metadata.CurrentPage,
		PageSize:     metadata.PageSize,
		FirstPage:    metadata.FirstPage,
		LastPage:     metadata.LastPage,
		TotalRecords: metadata.TotalRecords,
		Tasks:        tasks,
	}

	return tc, nil
}
