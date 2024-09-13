package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

type TaskLog struct {
	ID     uuid.UUID      `json:"id"`
	TaskID uuid.UUID      `json:"taskId"`
	Log    sql.NullString `json:"log"`
}

type TaskLogModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}

func (m *TaskLogModel) Get(ctx context.Context, id uuid.UUID) (*TaskLog, error) {
	query := `
SELECT id,
       task_id,
       log
FROM orchestrator.task_logs
WHERE id = $1;
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", database.MinifySQL(query)),
		slog.String("id", id.String()),
	))

	taskLog := &TaskLog{}
	logger.Info("performing query")
	err := m.Pool.QueryRow(ctx, query, id).Scan(
		&taskLog.ID,
		&taskLog.TaskID,
		&taskLog.Log,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.Error("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return taskLog, nil
}
