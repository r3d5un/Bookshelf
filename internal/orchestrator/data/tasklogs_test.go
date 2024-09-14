package data_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskLogModel(t *testing.T) {
	task := data.Task{
		Name:      "task_log_test",
		CronExpr:  sql.NullString{String: "* * * * *", Valid: true},
		Enabled:   sql.NullBool{Bool: false, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	_, err := models.Tasks.Insert(context.Background(), task)
	if err != nil {
		t.Errorf("error occurred while inserting task: %s\n", err)
		return
	}

	state := string(data.WaitingTaskState)
	timestamp := time.Now()
	taskQueueRow := data.TaskQueue{
		Name:      &task.Name,
		State:     &state,
		CreatedAt: &timestamp,
		UpdatedAt: &timestamp,
		RunAt:     &timestamp,
	}
	insertedTaskQueue, err := models.TaskQueues.Insert(context.Background(), taskQueueRow)
	if err != nil {
		t.Errorf("error occurred while inserting new task: %s", err)
		return
	}

	taskLog := data.TaskLog{
		ID:     uuid.New(),
		TaskID: insertedTaskQueue.ID,
		Log: string(
			`{"time":"2024-09-13T21:40:09.926641812+02:00","level":"INFO","msg":"test log statement"}`,
		),
	}

	tql := data.TaskLog{}

	t.Run("Insert", func(t *testing.T) {
		insertedTaskLog, err := models.TaskLogs.Insert(context.Background(), taskLog)
		if err != nil {
			t.Errorf("error occurred while inserting task log: %s\n", err)
			return
		}

		tql = *insertedTaskLog
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.TaskLogs.Get(context.Background(), tql.ID)
		if err != nil {
			t.Errorf("unable to get task log: %s\n", err)
			return
		}
	})

	t.Run("GetByTaskID", func(t *testing.T) {
		logs, err := models.TaskLogs.GetByTaskID(context.Background(), insertedTaskQueue.ID)
		if err != nil {
			t.Errorf("unablet o get task logs: %s\n", err)
			return
		}
		if len(logs) < 1 {
			t.Error("no logs returned")
			return
		}
	})
}
