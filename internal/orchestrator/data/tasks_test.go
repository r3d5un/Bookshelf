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
			t.Errorf("error occurred while querying task: %s", err)
			return
		}
	})
}
