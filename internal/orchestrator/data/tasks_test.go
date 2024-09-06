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
		Name:      sql.NullString{String: "test_task", Valid: true},
		CronExpr:  sql.NullString{String: "* * * * *", Valid: true},
		Enabled:   sql.NullBool{Bool: false, Valid: true},
		Deleted:   sql.NullBool{Bool: false, Valid: true},
		Timestamp: sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("Get", func(t *testing.T) {
		_, err := models.Tasks.Get(context.Background(), task.Name.String)
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
}
