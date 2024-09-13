package types

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type TaskLog struct {
	ID     uuid.UUID       `json:"id"`
	TaskID uuid.UUID       `json:"taskId"`
	Log    json.RawMessage `json:"log"`
}

func CreateTaskLog(ctx context.Context, models *data.Models, log TaskLog) (*TaskLog, error) {
	logRow := data.TaskLog{
		ID:     log.ID,
		TaskID: log.TaskID,
		Log:    log.Log,
	}

	insertedLogRow, err := models.TaskLogs.Insert(ctx, logRow)
	if err != nil {
		return nil, err
	}

	createdLog := TaskLog{
		ID:     insertedLogRow.ID,
		TaskID: insertedLogRow.TaskID,
		Log:    insertedLogRow.Log,
	}

	return &createdLog, nil
}
