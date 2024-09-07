package types

import (
	"context"
	"database/sql"
	"sync"
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
		Name:      taskRow.Name,
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
			Name:      t.Name,
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

func CreateTask(ctx context.Context, models *data.Models, task Task) (*Task, error) {
	dbRow := data.Task{
		Name:      task.Name,
		CronExpr:  newNullString(task.CronExpr),
		Enabled:   newNullBool(task.Enabled),
		UpdatedAt: newNullTime(task.UpdatedAt),
	}

	insertedTask, err := models.Tasks.Insert(ctx, dbRow)
	if err != nil {
		return nil, err
	}

	task = Task{
		Name:      insertedTask.Name,
		CronExpr:  nullStringToPtr(insertedTask.CronExpr),
		Enabled:   nullBoolToPtr(insertedTask.Enabled),
		UpdatedAt: nullTimeToPtr(insertedTask.UpdatedAt),
	}

	return &task, nil
}

func UpdateTask(ctx context.Context, models *data.Models, task Task) (*Task, error) {
	dbRow := data.Task{
		Name:      task.Name,
		CronExpr:  newNullString(task.CronExpr),
		Enabled:   newNullBool(task.Enabled),
		UpdatedAt: newNullTime(task.UpdatedAt),
	}

	updatedTask, err := models.Tasks.Update(ctx, dbRow)
	if err != nil {
		return nil, err
	}

	task = Task{
		Name:      updatedTask.Name,
		CronExpr:  nullStringToPtr(updatedTask.CronExpr),
		Enabled:   nullBoolToPtr(updatedTask.Enabled),
		UpdatedAt: nullTimeToPtr(updatedTask.UpdatedAt),
	}

	return &task, nil
}

func DeleteTask(ctx context.Context, models *data.Models, name string) (*Task, error) {
	deletedTask, err := models.Tasks.Delete(ctx, name)
	if err != nil {
		return nil, err
	}

	task := Task{
		Name:      deletedTask.Name,
		CronExpr:  nullStringToPtr(deletedTask.CronExpr),
		Enabled:   nullBoolToPtr(deletedTask.Enabled),
		UpdatedAt: nullTimeToPtr(deletedTask.UpdatedAt),
	}

	return &task, nil
}

func newNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func newNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func newNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func nullBoolToPtr(nb sql.NullBool) *bool {
	if !nb.Valid {
		return nil
	}
	return &nb.Bool
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
