package data_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskModel(t *testing.T) {
	task := data.Task{
		Name:      "test_task",
		CronExpr:  sql.NullString{String: "* * * * *", Valid: true},
		Enabled:   sql.NullBool{Bool: false, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("Insert", func(t *testing.T) {
		insertedTask, err := models.Tasks.Insert(context.Background(), task)
		if err != nil {
			t.Errorf("error occurred while inserting task: %s\n", err)
			return
		}

		if insertedTask.Name != task.Name {
			t.Errorf("expected task name %s, got %s\n", task.Name, insertedTask.Name)
			return
		}
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.Tasks.Get(context.Background(), task.Name)
		if err != nil {
			t.Errorf("error occurred while querying task: %s\n", err)
			return
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 100,
			OrderBy:  []string{"name"},
		}
		tasks, metadata, err := models.Tasks.GetAll(context.Background(), filters)
		if err != nil {
			t.Errorf("error occurred while reading tasks: %s\n", err)
			return
		}
		if len(tasks) < 1 {
			t.Errorf("no tasks returned")
			return
		}
		if metadata.CurrentPage != filters.Page {
			t.Errorf("expected page %d in metadata, got %d\n", filters.Page, metadata.CurrentPage)
			return
		}
	})

	t.Run("Update", func(t *testing.T) {
		task.CronExpr = sql.NullString{String: "1 * * * *", Valid: true}

		updatedTask, err := models.Tasks.Update(context.Background(), task)
		if err != nil {
			t.Errorf("error occurred while updating tasks: %s\n", err)
			return
		}
		if task.CronExpr.String != updatedTask.CronExpr.String {
			t.Errorf(
				"expected cron expression %s, got %s\n",
				task.CronExpr.String,
				updatedTask.CronExpr.String,
			)
		}
	})

	t.Run("UpsertOld", func(t *testing.T) {
		task.CronExpr = sql.NullString{String: "2 * * * *", Valid: true}

		upsertedTask, err := models.Tasks.Upsert(context.Background(), task)
		if err != nil {
			t.Errorf("error occurred while updating tasks: %s\n", err)
			return
		}
		if task.CronExpr.String != upsertedTask.CronExpr.String {
			t.Errorf(
				"expected cron expression %s, got %s\n",
				task.CronExpr.String,
				upsertedTask.CronExpr.String,
			)
		}
	})

	t.Run("UpsertNew", func(t *testing.T) {
		newTask := data.Task{
			Name:      "upserted_task",
			CronExpr:  sql.NullString{String: "* * * * *", Valid: true},
			Enabled:   sql.NullBool{Bool: false, Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		}

		upsertedTask, err := models.Tasks.Upsert(context.Background(), newTask)
		if err != nil {
			t.Errorf("error occurred while updating tasks: %s\n", err)
			return
		}
		if newTask.CronExpr.String != upsertedTask.CronExpr.String {
			t.Errorf(
				"expected cron expression %s, got %s\n",
				newTask.CronExpr.String,
				upsertedTask.CronExpr.String,
			)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := models.Tasks.Delete(context.Background(), task.Name)
		if err != nil {
			t.Errorf("an error occurred while deleting task: %s\n", err)
			return
		}
	})
}
