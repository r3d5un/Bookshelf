package types

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type Task struct {
	Name      string                      `json:"name"`
	CronExpr  *string                     `json:"cronExpr,omitempty"`
	Enabled   *bool                       `json:"enabled,omitempty"`
	UpdatedAt *time.Time                  `json:"timestamp,omitempty"`
	Job       func(context.Context) error `json:"-"`
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

// NewTask created a new task object. Not to be confused with CreateTask.
//
// NewTask does not persist data.
func NewTask(
	name string,
	cronExpr string,
	enabled bool,
	updatedAt time.Time,
	job func(context.Context) error,
) Task {
	return Task{
		Name:      name,
		CronExpr:  &cronExpr,
		Enabled:   &enabled,
		UpdatedAt: &updatedAt,
		Job:       job,
	}
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
	if filters.Page <= 0 {
		filters.Page = 1
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

// SyncTasks accepts a slice of tasks which from the caller, which acts
// as the master, syncing the task overview.
//
// New tasks are inserted in a disabled state.
//
// Old tasks are updated with data from the given list.
//
// Tasks in the task overview not found in the given task slice is deleted from
// the overview. The deletion is cascading, meaning all task run records will
// also be deleted.
func SyncTasks(ctx context.Context, models *data.Models, tasks []Task) error {
	logger := logging.LoggerFromContext(ctx)

	filters := data.Filters{
		PageSize: 50_000, // Set to a high value to retrieve all tasks
		OrderBy:  []string{"name"},
	}
	logger.Info("filters set", "filters", filters)

	logger.Info("reading existing tasks from database")
	tc, err := ReadAllTasks(ctx, models, filters)
	if err != nil {
		logger.Error("an error occurred while reading tasks", "error", err)
		return err
	}
	logger.Info("tasks read", "taskCollection", tc)

	logger.Info("checking for differences between wanted tasks and tasks in database")
	appTasks := make(map[string]Task, len(tasks))
	dbTasks := make(map[string]Task, len(tc.Tasks))

	for _, appTask := range tasks {
		appTasks[appTask.Name] = appTask
	}

	for _, dbTask := range tc.Tasks {
		dbTasks[dbTask.Name] = *dbTask
	}

	var newTasks []Task
	var deletableTasks []Task
	var updateableTasks []Task

	for _, appTask := range tasks {
		if _, found := dbTasks[appTask.Name]; found {
			// Tasks that already exists should not have their enabled status updated.
			// They should keep their enabled state as set by the users of the app,
			// independent of what the scheduler sets during their initialization.
			appTask.Enabled = dbTasks[appTask.Name].Enabled
			updateableTasks = append(updateableTasks, appTask)
		} else {
			enabled := false
			appTask.Enabled = &enabled
			newTasks = append(newTasks, appTask)
		}
	}

	for _, dbTask := range tc.Tasks {
		if _, found := appTasks[dbTask.Name]; !found {
			deletableTasks = append(deletableTasks, *dbTask)
		}
	}

	var wg sync.WaitGroup
	errorChan := make(chan error)

	logger.Info(
		"task differences set",
		"newTasks", newTasks,
		"deletableTasks", deletableTasks,
		"updateableTasks", updateableTasks,
	)

	logger.Info("updating database")
	for _, task := range newTasks {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := CreateTask(ctx, models, task)
			if err != nil {
				errorChan <- err
			}
		}()
	}

	for _, task := range deletableTasks {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := DeleteTask(ctx, models, task.Name)
			if err != nil {
				errorChan <- err
			}
		}()
	}

	for _, task := range updateableTasks {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := UpdateTask(ctx, models, task)
			if err != nil {
				errorChan <- err
			}
		}()
	}

	wg.Wait()
	close(errorChan)
	logger.Info("database updated")

	logger.Info("checking for errors from update process")
	for err := range errorChan {
		if err != nil {
			logger.Error("error occurred while syncing tasks with database", "error", err)
			return err
		}
	}

	logger.Info("database tasks updated")

	return nil
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
